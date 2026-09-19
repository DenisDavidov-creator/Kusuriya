package handler

import (
	"Kusuriya/internal/httpx"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type KnowledgebaseHandler struct {
	service KnowledgebaseServiceInterface
}

func NewKnowledgebaseHandler(service KnowledgebaseServiceInterface) *KnowledgebaseHandler {
	return &KnowledgebaseHandler{service: service}
}

// CreateCause godoc
// @Summary      Create cause
// @Description  Create cause
// @Tags         cause
// @Accept       json
// @Produce      json
// @Param        request body  CreateCauseRequest true "cause to create"
// @Success      200 {object}  CauseResponse
// @Failure      400 {object}  httpx.ErrorResponse
// @Failure      409 {object}  httpx.ErrorResponse
// @Failure 	 500 {object}  httpx.ErrorResponse
// @Router       /cause [post]
func (h *KnowledgebaseHandler) CreateCause(w http.ResponseWriter, r *http.Request) {
	var req CreateCauseRequest
	if err := httpx.Decode(r.Body, &req); err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}

	problems, err := req.Validate()
	if len(problems) != 0 {
		httpx.RespondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateCause(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCause(dom))
}

// GetCauseByID godoc
// @Summary      Get cause by id
// @Description  Returns one cause.
// @Tags         cause
// @Produce      json
// @Param        id   path      int  true  "cause id"
// @Success      200  {object}  CauseResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such cause"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause/{id} [get]
func (h *KnowledgebaseHandler) GetCauseByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetCauseByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCause(dom))
}

// GetCauses godoc
// @Summary      List causes
// @Description  Returns every cause, ordered by id.
// @Tags         cause
// @Produce      json
// @Success      200  {array}   CauseResponse
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause [get]
func (h *KnowledgebaseHandler) GetCauses(w http.ResponseWriter, r *http.Request) {
	doms, err := h.service.GetCauses(r.Context())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToCauses(doms))
}

// UpdateCause godoc
// @Summary      Update cause
// @Description  Replaces every field of a cause. The id in the body is ignored.
// @Tags         cause
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "cause id"
// @Param        request  body      UpdateCauseRequest  true  "new cause values"
// @Success      200      {object}  CauseResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      404  {object}  httpx.ErrorResponse  "no such cause"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause/{id} [put]
func (h *KnowledgebaseHandler) UpdateCause(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	var req UpdateCauseRequest
	if err := httpx.Decode(r.Body, &req); err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}

	problems, err := req.Validate()
	if len(problems) != 0 {
		httpx.RespondError(w, r, err, problems)
		return
	}

	dom, err := h.service.UpdateCause(r.Context(), req.ToDomain(id))
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCause(dom))
}

// DeleteCause godoc
// @Summary      Delete cause
// @Description  Deletes a cause and, by cascade, its symptom and factor profiles.
// @Tags         cause
// @Produce      json
// @Param        id   path      int  true  "cause id"
// @Success      202  {string}  string  "cause deleted"
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause/{id} [delete]
func (h *KnowledgebaseHandler) DeleteCause(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	err = h.service.DeleteCause(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusAccepted, "cause deleted")
}

// CreateRemedy godoc
// @Summary      Create remedy
// @Description  Adds a cure or antidote.
// @Tags         remedy
// @Accept       json
// @Produce      json
// @Param        request  body      CreateRemedyRequest  true  "create remedy"
// @Success      200      {object}  RemedyResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /remedy [post]
func (h *KnowledgebaseHandler) CreateRemedy(w http.ResponseWriter, r *http.Request) {
	var req CreateRemedyRequest
	if err := httpx.Decode(r.Body, &req); err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateRemedy(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToRemedy(dom))
}

// GetRemedyByID godoc
// @Summary      Get remedy by id
// @Description  Returns one remedy.
// @Tags         remedy
// @Produce      json
// @Param        id   path      int  true  "remedy id"
// @Success      200  {object}  RemedyResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such remedy"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /remedy/{id} [get]
func (h *KnowledgebaseHandler) GetRemedyByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetRemedyByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToRemedy(dom))
}

// GetRemedies godoc
// @Summary      List remedies
// @Description  Returns every remedy, ordered by id.
// @Tags         remedy
// @Produce      json
// @Success      200  {array}   RemedyResponse
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /remedy [get]
func (h *KnowledgebaseHandler) GetRemedies(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetRemedies(r.Context())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToRemedies(doms))
}

// CreateFactor godoc
// @Summary      Create factor
// @Description  Adds a contextual factor such as diet or season.
// @Tags         factor
// @Accept       json
// @Produce      json
// @Param        request  body      CreateFactorRequest  true  "create factor"
// @Success      200      {object}  FactorResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /factor [post]
func (h *KnowledgebaseHandler) CreateFactor(w http.ResponseWriter, r *http.Request) {
	var req CreateFactorRequest
	if err := httpx.Decode(r.Body, &req); err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problems)
		return
	}

	dom, err := h.service.CreateFactor(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToFactor(dom))
}

// GetFactorByID godoc
// @Summary      Get factor by id
// @Description  Returns one factor.
// @Tags         factor
// @Produce      json
// @Param        id   path      int  true  "factor id"
// @Success      200  {object}  FactorResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such factor"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /factor/{id} [get]
func (h *KnowledgebaseHandler) GetFactorByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetFactorByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToFactor(dom))
}

// GetFactors godoc
// @Summary      List factors
// @Description  Returns every factor, ordered by id.
// @Tags         factor
// @Produce      json
// @Success      200  {array}   FactorResponse
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /factor [get]
func (h *KnowledgebaseHandler) GetFactors(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetFactors(r.Context())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToFactors(doms))
}

// CreateCategory godoc
// @Summary      Create category
// @Description  Adds a symptom category.
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCategoryRequest  true  "create category"
// @Success      200      {object}  CategoryResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /category [post]
func (h *KnowledgebaseHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	err := httpx.Decode(r.Body, &req)

	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}

	problmes, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problmes)
		return
	}

	dom, err := h.service.CreateCategory(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCategory(dom))
}

// GetCategoryByID godoc
// @Summary      Get category by id
// @Description  Returns one category.
// @Tags         category
// @Produce      json
// @Param        id   path      int  true  "category id"
// @Success      200  {object}  CategoryResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such category"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /category/{id} [get]
func (h *KnowledgebaseHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToCategory(dom))
}

// GetCategories godoc
// @Summary      List categories
// @Description  Returns every category, ordered by id.
// @Tags         category
// @Produce      json
// @Success      200  {array}   CategoryResponse
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /category [get]
func (h *KnowledgebaseHandler) GetCategories(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetCategories(r.Context())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToCategories(doms))
}

// CreateSymptom godoc
// @Summary      Create symptom
// @Description  Adds a symptom belonging to a category.
// @Tags         symptom
// @Accept       json
// @Produce      json
// @Param        request  body      CreateSymptomRequest  true  "create symptom"
// @Success      200      {object}  SymptomResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /symptom [post]
func (h *KnowledgebaseHandler) CreateSymptom(w http.ResponseWriter, r *http.Request) {
	var req CreateSymptomRequest
	err := httpx.Decode(r.Body, &req)

	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}

	problmes, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problmes)
		return
	}

	dom, err := h.service.CreateSymptom(r.Context(), req.ToDomain())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToSymptom(dom))
}

// GetSymptomByID godoc
// @Summary      Get symptom by id
// @Description  Returns one symptom.
// @Tags         symptom
// @Produce      json
// @Param        id   path      int  true  "symptom id"
// @Success      200  {object}  SymptomResponse
// @Failure      400  {object}  httpx.ErrorResponse  "id is not a number"
// @Failure      404  {object}  httpx.ErrorResponse  "no such symptom"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /symptom/{id} [get]
func (h *KnowledgebaseHandler) GetSymptomByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrong URLParam")
		return
	}

	dom, err := h.service.GetSymptomByID(r.Context(), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainToSymptom(dom))
}

// GetSymptoms godoc
// @Summary      List symptoms
// @Description  Returns every symptom, ordered by id.
// @Tags         symptom
// @Produce      json
// @Success      200  {array}   SymptomResponse
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /symptom [get]
func (h *KnowledgebaseHandler) GetSymptoms(w http.ResponseWriter, r *http.Request) {

	doms, err := h.service.GetSymptoms(r.Context())
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToSymptoms(doms))
}

// CreateSymptomProfiles godoc
// @Summary      Set symptom profiles of a cause
// @Description  Adds P(symptom | cause) entries for one cause.
// @Tags         cause
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "cause id"
// @Param        request  body      CreateSymptomsProfileRequests  true  "profiles to store"
// @Success      200      {array}   SymptomProfileResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause/{id}/symptom-profiles [put]
func (h *KnowledgebaseHandler) CreateSymptomProfiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrongURLParam")
		return
	}
	var req CreateSymptomsProfileRequests
	err = httpx.Decode(r.Body, &req)
	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problems)
		return
	}
	doms, err := h.service.CreateSymptomProfiles(r.Context(), req.ToDomains(int64(id)), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToSymptomsProfiles(doms))
}

// CreateFactorProfiles godoc
// @Summary      Set factor profiles of a cause
// @Description  Adds P(factor | cause) modifiers for one cause.
// @Tags         cause
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "cause id"
// @Param        request  body      CreateFactorsProfileRequests  true  "profiles to store"
// @Success      200      {array}   FactorProfileResponse
// @Failure      400  {object}  httpx.ErrorResponse  "malformed JSON or validation failed"
// @Failure      409  {object}  httpx.ErrorResponse  "code already exists"
// @Failure      500  {object}  httpx.ErrorResponse  "unexpected error"
// @Router       /cause/{id}/factor-profiles [put]
func (h *KnowledgebaseHandler) CreateFactorProfiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondError(w, r, err, "wrongURLParam")
		return
	}
	var req CreateFactorsProfileRequests
	err = httpx.Decode(r.Body, &req)
	if err != nil {
		httpx.RespondError(w, r, err, err.Error())
		return
	}
	problems, err := req.Validate()
	if err != nil {
		httpx.RespondError(w, r, err, problems)
		return
	}
	doms, err := h.service.CreateFactorProfiles(r.Context(), req.ToDomains(int64(id)), id)
	if err != nil {
		httpx.RespondError(w, r, err, "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, DomainsToFactorsProfiles(doms))
}
