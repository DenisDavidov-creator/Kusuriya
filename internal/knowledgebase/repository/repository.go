package repository

import (
	"Kusuriya/internal/domain"
	sqlcgen "Kusuriya/internal/knowledgebase/repository/sqlc"
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KnowledgebaseRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewKnowledgebaseRepository(pool *pgxpool.Pool) *KnowledgebaseRepository {
	return &KnowledgebaseRepository{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}
}

type transactionalKey string

var trKey transactionalKey = "TrKey"

func (r *KnowledgebaseRepository) WithTx(ctx context.Context) (context.Context, error) {
	trn, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	ctxWV := context.WithValue(ctx, trKey, trn)
	return ctxWV, nil
}

func (r *KnowledgebaseRepository) Commit(ctx context.Context) error {
	trn, ok := ctx.Value(trKey).(pgx.Tx)
	if ok {
		err := trn.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrTransaction
}

func (r *KnowledgebaseRepository) Rollback(ctx context.Context) error {
	trn, ok := ctx.Value(trKey).(pgx.Tx)
	if ok {
		err := trn.Rollback(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrTransaction
}

func (r *KnowledgebaseRepository) txOrNot(ctx context.Context) *sqlcgen.Queries {
	var pool *sqlcgen.Queries
	tx, ok := ctx.Value(trKey).(pgx.Tx)
	if ok {
		pool = r.queries.WithTx(tx)
	} else {
		pool = r.queries
	}
	return pool
}

func (r *KnowledgebaseRepository) CreateCause(ctx context.Context, d domain.Cause) (*domain.Cause, error) {
	var remedyID pgtype.Int8
	if d.RemedyID != nil {
		remedyID = pgtype.Int8{Int64: *d.RemedyID, Valid: true}
	}
	params := sqlcgen.CreateCauseParams{
		Code:  d.Code,
		Name:  d.Name,
		Type:  d.Type,
		Prior: d.Prior,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		RemedyID: remedyID,
	}

	result, err := r.queries.CreateCause(ctx, params)

	if err != nil {
		return nil, MapPgError(err)
	}

	dom := &domain.Cause{
		ID:    result.ID,
		Code:  result.Code,
		Name:  result.Name,
		Type:  result.Type,
		Prior: result.Prior,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.RemedyID.Valid {
		id := result.RemedyID.Int64
		dom.RemedyID = &id
	}

	return dom, nil
}

func (r *KnowledgebaseRepository) GetCauseByID(ctx context.Context, id int) (*domain.Cause, error) {
	result, err := r.queries.GetCauseByID(ctx, int64(id))
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Cause{
		ID:    result.ID,
		Code:  result.Code,
		Name:  result.Name,
		Type:  result.Type,
		Prior: result.Prior,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.RemedyID.Valid {
		id := result.RemedyID.Int64
		dom.RemedyID = &id
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetCauses(ctx context.Context) ([]domain.Cause, error) {
	results, err := r.queries.GetCauses(ctx)
	if err != nil {
		return nil, MapPgError(err)
	}
	doms := make([]domain.Cause, 0, len(results))

	for _, result := range results {
		dom := domain.Cause{
			ID:    result.ID,
			Code:  result.Code,
			Name:  result.Name,
			Type:  result.Type,
			Prior: result.Prior,
		}
		if result.Description.Valid {
			dom.Description = result.Description.String
		}
		if result.RemedyID.Valid {
			id := result.RemedyID.Int64
			dom.RemedyID = &id
		}
		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) UpdateCause(ctx context.Context, d domain.Cause) (*domain.Cause, error) {
	var remedyID pgtype.Int8
	if d.RemedyID != nil {
		remedyID = pgtype.Int8{Int64: *d.RemedyID, Valid: true}
	}
	params := sqlcgen.UpdateCauseParams{
		Code:  d.Code,
		Name:  d.Name,
		Type:  d.Type,
		Prior: d.Prior,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		RemedyID: remedyID,
	}

	result, err := r.queries.UpdateCause(ctx, params)
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Cause{
		ID:    result.ID,
		Code:  result.Code,
		Name:  result.Name,
		Type:  result.Type,
		Prior: result.Prior,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.RemedyID.Valid {
		id := result.RemedyID.Int64
		dom.RemedyID = &id
	}

	return dom, nil
}

func (r *KnowledgebaseRepository) DeleteCause(ctx context.Context, id int) error {
	err := r.queries.DeleteCause(ctx, int64(id))
	if err != nil {
		return MapPgError(err)
	}
	return nil
}

func (r *KnowledgebaseRepository) CreateRemedy(ctx context.Context, d domain.Remedy) (*domain.Remedy, error) {
	params := sqlcgen.CreateRemedyParams{
		Name: d.Name,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		Notes: pgtype.Text{
			String: d.Notes,
			Valid:  d.Notes != "",
		},
	}

	result, err := r.queries.CreateRemedy(ctx, params)
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Remedy{
		ID:   result.ID,
		Name: result.Name,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.Notes.Valid {
		dom.Notes = result.Notes.String
	}

	return dom, nil
}

func (r *KnowledgebaseRepository) GetCausesWithRemediesByIDs(ctx context.Context, IDs []int64) ([]domain.CauseWithRemedy, error) {
	results, err := r.queries.GetCausesWithRemediesByIDs(ctx, IDs)
	if err != nil {
		return nil, MapPgError(err)
	}
	var doms = make([]domain.CauseWithRemedy, 0, len(results))
	for _, result := range results {
		dom := domain.CauseWithRemedy{
			ID:    result.ID,
			Code:  result.Code,
			Name:  result.Name,
			Type:  result.Type,
			Prior: result.Prior,
		}
		if result.Description.Valid {
			dom.Description = result.Description.String
		}
		if result.RemedyID.Valid {
			id := result.RemedyID.Int64
			dom.RemedyID = &id
		}
		if result.RemedyName.Valid {
			remedyName := result.RemedyName.String
			dom.RemedyName = &remedyName
		}
		if result.RemedyDescription.Valid {
			remedyDescription := result.RemedyDescription.String
			dom.RemedyDescription = &remedyDescription
		}
		if result.RemedyNotes.Valid {
			remedyNotes := result.RemedyNotes.String
			dom.RemedyNotes = &remedyNotes
		}
		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) GetRemedyByID(ctx context.Context, id int) (*domain.Remedy, error) {
	result, err := r.queries.GetRemedyByID(ctx, int64(id))
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Remedy{
		ID:   result.ID,
		Name: result.Name,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	if result.Notes.Valid {
		dom.Notes = result.Notes.String
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetRemedies(ctx context.Context) ([]domain.Remedy, error) {
	result, err := r.queries.GetRemedies(ctx)
	if err != nil {
		return nil, MapPgError(err)
	}

	doms := make([]domain.Remedy, 0, len(result))
	for _, result := range result {
		dom := domain.Remedy{
			ID:   result.ID,
			Name: result.Name,
		}
		if result.Description.Valid {
			dom.Description = result.Description.String
		}
		if result.Notes.Valid {
			dom.Notes = result.Notes.String
		}
		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) CreateFactor(ctx context.Context, d domain.Factor) (*domain.Factor, error) {
	params := sqlcgen.CreateFactorParams{
		Code: d.Code,
		Name: d.Name,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
	}

	result, err := r.queries.CreateFactor(ctx, params)
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Factor{
		ID:   result.ID,
		Code: result.Code,
		Name: result.Name,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetFactorByID(ctx context.Context, id int) (*domain.Factor, error) {
	result, err := r.queries.GetFactorByID(ctx, int64(id))
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Factor{
		ID:   result.ID,
		Name: result.Name,
		Code: result.Code,
	}
	if result.Description.Valid {
		dom.Description = result.Description.String
	}

	return dom, nil
}

func (r *KnowledgebaseRepository) GetFactors(ctx context.Context) ([]domain.Factor, error) {
	result, err := r.queries.GetFactors(ctx)
	if err != nil {
		return nil, MapPgError(err)
	}

	doms := make([]domain.Factor, 0, len(result))
	for _, result := range result {
		dom := domain.Factor{
			ID:   result.ID,
			Name: result.Name,
			Code: result.Code,
		}
		if result.Description.Valid {
			dom.Description = result.Description.String
		}

		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) CreateCategory(ctx context.Context, d domain.Category) (*domain.Category, error) {

	result, err := r.queries.CreateCategory(ctx, d.Name)
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Category{
		ID:   result.ID,
		Name: result.Name,
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	result, err := r.queries.GetCategoryByID(ctx, int64(id))
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Category{
		ID:   result.ID,
		Name: result.Name,
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	result, err := r.queries.GetCategories(ctx)
	if err != nil {
		return nil, MapPgError(err)
	}
	doms := make([]domain.Category, 0, len(result))
	for _, result := range result {
		dom := domain.Category{
			ID:   result.ID,
			Name: result.Name,
		}

		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) CreateSymptom(ctx context.Context, d domain.Symptom) (*domain.Symptom, error) {

	params := sqlcgen.CreateSymptomParams{
		Code: d.Code,
		Name: d.Name,
		Description: pgtype.Text{
			String: d.Description,
			Valid:  d.Description != "",
		},
		CategoryID: d.CategoryID,
	}

	result, err := r.queries.CreateSymptom(ctx, params)
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Symptom{
		ID:         result.ID,
		Code:       result.Code,
		Name:       result.Name,
		CategoryID: result.CategoryID,
	}
	if result.Description.String != "" {
		dom.Description = result.Description.String
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetSymptomByID(ctx context.Context, id int) (*domain.Symptom, error) {
	result, err := r.queries.GetSymptomByID(ctx, int64(id))
	if err != nil {
		return nil, MapPgError(err)
	}
	dom := &domain.Symptom{
		ID:         result.ID,
		Code:       result.Code,
		Name:       result.Name,
		CategoryID: result.CategoryID,
	}
	if result.Description.String != "" {
		dom.Description = result.Description.String
	}
	return dom, nil
}

func (r *KnowledgebaseRepository) GetSymptoms(ctx context.Context) ([]domain.Symptom, error) {
	result, err := r.queries.GetSymptoms(ctx)
	if err != nil {
		return nil, MapPgError(err)
	}
	doms := make([]domain.Symptom, 0, len(result))
	for _, result := range result {
		dom := domain.Symptom{
			ID:         result.ID,
			Code:       result.Code,
			Name:       result.Name,
			CategoryID: result.CategoryID,
		}
		if result.Description.String != "" {
			dom.Description = result.Description.String
		}

		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) CreateSymptomProfiles(ctx context.Context, ds []domain.SymptomProfile) ([]domain.SymptomProfile, error) {
	pool := r.txOrNot(ctx)
	doms := make([]domain.SymptomProfile, 0, len(ds))
	for _, d := range ds {
		likelihood := pgtype.Numeric{}
		like := strconv.FormatFloat(d.Likelihood, 'f', -1, 64)
		err := likelihood.Scan(like)
		if err != nil {
			return nil, MapPgError(err)
		}

		params := sqlcgen.CreateSymptomProfilesParams{
			CauseID:    d.CauseID,
			SymptomID:  d.SymptomID,
			Likelihood: likelihood,
		}
		result, err := pool.CreateSymptomProfiles(ctx, params)
		if err != nil {
			return nil, MapPgError(err)
		}
		dom := domain.SymptomProfile{
			ID:        result.ID,
			CauseID:   result.CauseID,
			SymptomID: result.SymptomID,
		}
		l, err := result.Likelihood.Float64Value()
		if err != nil {
			return nil, MapPgError(err)
		}
		if l.Valid {
			dom.Likelihood = l.Float64
		}
		doms = append(doms, dom)
	}

	return doms, nil
}

func (r *KnowledgebaseRepository) DeleteSymptomProfiles(ctx context.Context, id int64) error {
	pool := r.txOrNot(ctx)
	err := pool.DeleteSymptomProfiles(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *KnowledgebaseRepository) CreateFactorProfiles(ctx context.Context, ds []domain.FactorProfile) ([]domain.FactorProfile, error) {
	pool := r.txOrNot(ctx)
	doms := make([]domain.FactorProfile, 0, len(ds))
	for _, d := range ds {
		modifier := pgtype.Numeric{}
		mod := strconv.FormatFloat(d.Modifier, 'f', -1, 64)
		err := modifier.Scan(mod)
		if err != nil {
			return nil, MapPgError(err)
		}

		params := sqlcgen.CreateFactorProfilesParams{
			CauseID:  d.CauseID,
			FactorID: d.FactorID,
			Modifier: modifier,
		}
		result, err := pool.CreateFactorProfiles(ctx, params)
		if err != nil {
			return nil, MapPgError(err)
		}
		dom := domain.FactorProfile{
			ID:       result.ID,
			CauseID:  result.CauseID,
			FactorID: result.FactorID,
		}
		l, err := result.Modifier.Float64Value()
		if err != nil {
			return nil, MapPgError(err)
		}
		if l.Valid {
			dom.Modifier = l.Float64
		}
		doms = append(doms, dom)
	}

	return doms, nil
}

func (r *KnowledgebaseRepository) DeleteFactorProfiles(ctx context.Context, id int64) error {
	pool := r.txOrNot(ctx)
	err := pool.DeleteFactorProfiles(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *KnowledgebaseRepository) GetFactorsByCodes(ctx context.Context, codes []string) ([]domain.FactorEvidence, error) {
	results, err := r.queries.GetFactorsByCodes(ctx, codes)
	if err != nil {
		return nil, err
	}
	var doms []domain.FactorEvidence
	for _, result := range results {
		dom := domain.FactorEvidence{
			CauseID: result.ID,
		}
		l, err := result.Modifier.Float64Value()
		if err != nil {
			return nil, MapPgError(err)
		}
		if l.Valid {
			m := l.Float64
			dom.Modifier = &m
		} else {
			dom.Modifier = nil
		}
		doms = append(doms, dom)
	}
	return doms, nil
}

func (r *KnowledgebaseRepository) GetSymptomsByCodes(ctx context.Context, codes []string) ([]domain.SymptomEvidence, error) {
	results, err := r.queries.GetSymptomsByCodes(ctx, codes)
	if err != nil {
		return nil, err
	}
	var doms []domain.SymptomEvidence
	for _, result := range results {
		dom := domain.SymptomEvidence{
			CauseID: result.ID,
		}
		l, err := result.Likelihood.Float64Value()
		if err != nil {
			return nil, MapPgError(err)
		}
		if l.Valid {
			lh := l.Float64
			dom.Likelihood = &lh
		} else {
			dom.Likelihood = nil
		}
		doms = append(doms, dom)
	}
	return doms, nil
}
