package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/RafayLabs/rcloud-base/internal/dao"
	"github.com/RafayLabs/rcloud-base/internal/models"
	commonv3 "github.com/RafayLabs/rcloud-base/proto/types/commonpb/v3"
	systemv3 "github.com/RafayLabs/rcloud-base/proto/types/systempb/v3"
	"github.com/google/uuid"
	bun "github.com/uptrace/bun"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type OIDCProviderService interface {
	Create(context.Context, *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error)
	GetByID(context.Context, *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error)
	GetByName(context.Context, *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error)
	List(context.Context) (*systemv3.OIDCProviderList, error)
	Update(context.Context, *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error)
	Delete(context.Context, *systemv3.OIDCProvider) error
}

type oidcProvider struct {
	db        *bun.DB
	kratosUrl string
	al        *zap.Logger
}

func NewOIDCProviderService(db *bun.DB, kratosUrl string, al *zap.Logger) OIDCProviderService {
	return &oidcProvider{db: db, kratosUrl: kratosUrl, al: al}
}

func generateCallbackUrl(id string, kUrl string) string {
	b, _ := url.Parse(kUrl)
	return fmt.Sprintf("%s://%s/self-service/methods/oidc/callback/%s", b.Scheme, b.Host, id)
}

func validateURL(rawURL string) error {
	_, err := url.ParseRequestURI(rawURL)
	return err
}

func (s *oidcProvider) getPartnerOrganization(ctx context.Context, provider *systemv3.OIDCProvider) (uuid.UUID, uuid.UUID, error) {
	partner := provider.GetMetadata().GetPartner()
	org := provider.GetMetadata().GetOrganization()
	partnerId, err := dao.GetPartnerId(ctx, s.db, partner)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	organizationId, err := dao.GetOrganizationId(ctx, s.db, org)
	if err != nil {
		return partnerId, uuid.Nil, err
	}
	return partnerId, organizationId, nil
}

func (s *oidcProvider) Create(ctx context.Context, provider *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error) {
	name := provider.Metadata.GetName()
	if len(name) == 0 {
		return &systemv3.OIDCProvider{}, fmt.Errorf("EMPTY NAME")
	}
	scopes := provider.GetSpec().GetScopes()
	if scopes == nil || len(scopes) == 0 {
		return &systemv3.OIDCProvider{}, fmt.Errorf("EMPTY SCOPES")
	}

	partnerId, organizationId, err := s.getPartnerOrganization(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("unable to get partner and org id")
	}
	p, _ := dao.GetIdByNamePartnerOrg(
		ctx,
		s.db,
		provider.GetMetadata().GetName(),
		uuid.NullUUID{UUID: partnerId, Valid: true},
		uuid.NullUUID{UUID: organizationId, Valid: true},
		&models.OIDCProvider{},
	)
	if p != nil {
		return nil, fmt.Errorf("OIDC provider %q already exists", provider.GetMetadata().GetName())
	}

	mapUrl := provider.Spec.GetMapperUrl()
	issUrl := provider.Spec.GetIssuerUrl()
	authUrl := provider.Spec.GetAuthUrl()
	tknUrl := provider.Spec.GetTokenUrl()

	if len(mapUrl) != 0 && validateURL(mapUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID MAPPER URL")
	}
	if len(issUrl) != 0 && validateURL(issUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID ISSUER URL")
	}
	if len(authUrl) != 0 && validateURL(authUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID AUTH URL")
	}
	if len(tknUrl) != 0 && validateURL(tknUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID TOKEN URL")
	}

	entity := &models.OIDCProvider{
		Name:            name,
		Description:     provider.GetMetadata().GetDescription(),
		CreatedAt:       time.Time{},
		ModifiedAt:      time.Time{},
		PartnerId:       partnerId,
		OrganizationId:  organizationId,
		ProviderName:    provider.Spec.GetProviderName(),
		MapperURL:       mapUrl,
		MapperFilename:  provider.Spec.GetMapperFilename(),
		ClientId:        provider.Spec.GetClientId(),
		ClientSecret:    provider.Spec.GetClientSecret(),
		Scopes:          provider.Spec.GetScopes(),
		IssuerURL:       issUrl,
		AuthURL:         authUrl,
		TokenURL:        tknUrl,
		RequestedClaims: provider.Spec.GetRequestedClaims().AsMap(),
		Predefined:      provider.Spec.GetPredefined(),
	}
	_, err = dao.Create(ctx, s.db, entity)
	if err != nil {
		return &systemv3.OIDCProvider{}, err
	}

	rclaims, _ := structpb.NewStruct(entity.RequestedClaims)
	rv := &systemv3.OIDCProvider{
		ApiVersion: apiVersion,
		Kind:       "OIDCProvider",
		Metadata: &commonv3.Metadata{
			Name:        entity.Name,
			Description: entity.Description,
			Id:          entity.Id.String(),
		},
		Spec: &systemv3.OIDCProviderSpec{
			ProviderName:    entity.ProviderName,
			MapperUrl:       entity.MapperURL,
			MapperFilename:  entity.MapperFilename,
			ClientId:        entity.ClientId,
			Scopes:          entity.Scopes,
			IssuerUrl:       entity.IssuerURL,
			AuthUrl:         entity.AuthURL,
			TokenUrl:        entity.TokenURL,
			RequestedClaims: rclaims,
			Predefined:      entity.Predefined,
			CallbackUrl:     generateCallbackUrl(entity.Id.String(), s.kratosUrl),
		},
	}

	CreateOidcAuditEvent(ctx, s.al, AuditActionCreate, rv.GetMetadata().GetName(), entity.Id)
	return rv, nil
}

func (s *oidcProvider) GetByID(ctx context.Context, provider *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error) {
	id, err := uuid.Parse(provider.Metadata.GetId())
	if err != nil {
		return &systemv3.OIDCProvider{}, err
	}

	entity := &models.OIDCProvider{}
	_, err = dao.GetByID(ctx, s.db, id, entity)
	// TODO: Return proper error for Id not exist
	if err != nil {
		return &systemv3.OIDCProvider{}, err
	}

	rclaims, _ := structpb.NewStruct(entity.RequestedClaims)
	rv := &systemv3.OIDCProvider{
		ApiVersion: apiVersion,
		Kind:       "OIDCProvider",
		Metadata: &commonv3.Metadata{
			Name:        entity.Name,
			Description: entity.Description,
			Id:          entity.Id.String(),
		},
		Spec: &systemv3.OIDCProviderSpec{
			ProviderName:    entity.ProviderName,
			MapperUrl:       entity.MapperURL,
			MapperFilename:  entity.MapperFilename,
			ClientId:        entity.ClientId,
			Scopes:          entity.Scopes,
			IssuerUrl:       entity.IssuerURL,
			AuthUrl:         entity.AuthURL,
			TokenUrl:        entity.TokenURL,
			RequestedClaims: rclaims,
			Predefined:      entity.Predefined,
			CallbackUrl:     generateCallbackUrl(entity.Id.String(), s.kratosUrl),
		},
	}
	return rv, nil
}

func (s *oidcProvider) GetByName(ctx context.Context, provider *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error) {
	name := provider.Metadata.GetName()
	if len(name) == 0 {
		return &systemv3.OIDCProvider{}, status.Error(codes.InvalidArgument, "EMPTY NAME")
	}

	entity := &models.OIDCProvider{}
	_, err := dao.GetByName(ctx, s.db, name, entity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &systemv3.OIDCProvider{}, status.Errorf(codes.InvalidArgument, "OIDC PROVIDER %q NOT EXIST", name)
		} else {
			return &systemv3.OIDCProvider{}, status.Errorf(codes.Internal, codes.Internal.String())
		}

	}

	rclaims, _ := structpb.NewStruct(entity.RequestedClaims)
	rv := &systemv3.OIDCProvider{
		ApiVersion: apiVersion,
		Kind:       "OIDCProvider",
		Metadata: &commonv3.Metadata{
			Name:         entity.Name,
			Description:  entity.Description,
			Id:           entity.Id.String(),
			Organization: entity.OrganizationId.String(),
			Partner:      entity.PartnerId.String(),
		},
		Spec: &systemv3.OIDCProviderSpec{
			ProviderName:    entity.ProviderName,
			MapperUrl:       entity.MapperURL,
			MapperFilename:  entity.MapperFilename,
			ClientId:        entity.ClientId,
			Scopes:          entity.Scopes,
			IssuerUrl:       entity.IssuerURL,
			AuthUrl:         entity.AuthURL,
			TokenUrl:        entity.TokenURL,
			RequestedClaims: rclaims,
			Predefined:      entity.Predefined,
			CallbackUrl:     generateCallbackUrl(entity.Id.String(), s.kratosUrl),
		},
	}
	return rv, nil
}

func (s *oidcProvider) List(ctx context.Context) (*systemv3.OIDCProviderList, error) {
	var (
		entities []models.OIDCProvider
		orgID    uuid.NullUUID
		parID    uuid.NullUUID
	)
	_, err := dao.List(ctx, s.db, parID, orgID, &entities)
	if err != nil {
		return &systemv3.OIDCProviderList{}, nil
	}
	var result []*systemv3.OIDCProvider
	for _, entity := range entities {
		rclaims, _ := structpb.NewStruct(entity.RequestedClaims)
		e := &systemv3.OIDCProvider{
			ApiVersion: apiVersion,
			Kind:       "OIDCProvider",
			Metadata: &commonv3.Metadata{
				Name:        entity.Name,
				Description: entity.Description,
				Id:          entity.Id.String(),
			},
			Spec: &systemv3.OIDCProviderSpec{
				ProviderName:    entity.ProviderName,
				MapperUrl:       entity.MapperURL,
				MapperFilename:  entity.MapperFilename,
				ClientId:        entity.ClientId,
				Scopes:          entity.Scopes,
				IssuerUrl:       entity.IssuerURL,
				AuthUrl:         entity.AuthURL,
				TokenUrl:        entity.TokenURL,
				RequestedClaims: rclaims,
				Predefined:      entity.Predefined,
				CallbackUrl:     generateCallbackUrl(entity.Id.String(), s.kratosUrl),
			},
		}
		result = append(result, e)
	}

	rv := &systemv3.OIDCProviderList{
		ApiVersion: "usermgmt.k8smgmt.io/v3",
		Kind:       "OIDCProviderList",
		Items:      result,
	}
	return rv, nil
}

func (s *oidcProvider) Update(ctx context.Context, provider *systemv3.OIDCProvider) (*systemv3.OIDCProvider, error) {
	name := provider.GetMetadata().GetName()
	if len(name) == 0 {
		return &systemv3.OIDCProvider{}, status.Error(codes.InvalidArgument, "EMPTY NAME")
	}
	scopes := provider.GetSpec().GetScopes()
	if scopes == nil || len(scopes) == 0 {
		return &systemv3.OIDCProvider{}, fmt.Errorf("EMPTY SCOPES")
	}

	partnerId, organizationId, err := s.getPartnerOrganization(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("unable to get partner and org id")
	}

	existingP := &models.OIDCProvider{}
	_, err = dao.GetByName(ctx, s.db, name, existingP)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &systemv3.OIDCProvider{}, status.Errorf(codes.InvalidArgument, "OIDC PROVIDER %q NOT EXIST", name)
		} else {
			return &systemv3.OIDCProvider{}, status.Error(codes.Internal, codes.Internal.String())
		}
	}

	mapUrl := provider.Spec.GetMapperUrl()
	issUrl := provider.Spec.GetIssuerUrl()
	authUrl := provider.Spec.GetAuthUrl()
	tknUrl := provider.Spec.GetTokenUrl()

	if len(mapUrl) != 0 && validateURL(mapUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID MAPPER URL")
	}
	if len(issUrl) != 0 && validateURL(issUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID ISSUER URL")
	}
	if len(authUrl) != 0 && validateURL(authUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID AUTH URL")
	}
	if len(tknUrl) != 0 && validateURL(tknUrl) != nil {
		return &systemv3.OIDCProvider{}, fmt.Errorf("INVALID TOKEN URL")
	}

	entity := &models.OIDCProvider{
		Name:            provider.Metadata.GetName(),
		Description:     provider.Metadata.GetDescription(),
		OrganizationId:  organizationId,
		PartnerId:       partnerId,
		ModifiedAt:      time.Now(),
		ProviderName:    provider.Spec.GetProviderName(),
		MapperURL:       mapUrl,
		MapperFilename:  provider.Spec.GetMapperFilename(),
		ClientId:        provider.Spec.GetClientId(),
		Scopes:          provider.Spec.GetScopes(),
		IssuerURL:       issUrl,
		AuthURL:         authUrl,
		TokenURL:        tknUrl,
		RequestedClaims: provider.Spec.GetRequestedClaims().AsMap(),
		Predefined:      provider.Spec.GetPredefined(),
	}
	_, err = dao.Update(ctx, s.db, existingP.Id, entity)
	if err != nil {
		return &systemv3.OIDCProvider{}, err
	}

	rclaims, _ := structpb.NewStruct(entity.RequestedClaims)
	rv := &systemv3.OIDCProvider{
		ApiVersion: apiVersion,
		Kind:       "OIDCProvider",
		Metadata: &commonv3.Metadata{
			Name:        entity.Name,
			Description: entity.Description,
			Id:          provider.GetMetadata().GetId(),
		},
		Spec: &systemv3.OIDCProviderSpec{
			ProviderName:    entity.ProviderName,
			MapperUrl:       entity.MapperURL,
			MapperFilename:  entity.MapperFilename,
			ClientId:        entity.ClientId,
			Scopes:          entity.Scopes,
			IssuerUrl:       entity.IssuerURL,
			AuthUrl:         entity.AuthURL,
			TokenUrl:        entity.TokenURL,
			RequestedClaims: rclaims,
			Predefined:      entity.Predefined,
			CallbackUrl:     generateCallbackUrl(provider.GetMetadata().GetId(), s.kratosUrl),
		},
	}

	CreateOidcAuditEvent(ctx, s.al, AuditActionUpdate, rv.GetMetadata().GetName(), entity.Id)
	return rv, nil
}

func (s *oidcProvider) Delete(ctx context.Context, provider *systemv3.OIDCProvider) error {
	entity := &models.OIDCProvider{}
	name := provider.GetMetadata().GetName()
	if len(name) == 0 {
		return status.Error(codes.InvalidArgument, "EMPTY NAME")
	}
	_, err := dao.GetByName(ctx, s.db, name, entity)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "OIDC PROVIDER %q NOT EXIST", name)
	}

	err = dao.Delete(ctx, s.db, entity.Id, &models.OIDCProvider{})
	if err != nil {
		return err
	}

	CreateOidcAuditEvent(ctx, s.al, AuditActionDelete, name, entity.Id)
	return nil
}