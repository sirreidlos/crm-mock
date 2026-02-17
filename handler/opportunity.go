package handler

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/go-chi/chi/v5"
    "crm-mock/model"
)

type OpportunityHandler struct { db *sql.DB }

func NewOpportunityHandler(db *sql.DB) *OpportunityHandler {
    return &OpportunityHandler{db: db}
}

func (h *OpportunityHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.CreateOpportunityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    id := fmt.Sprintf("OPP-%d", time.Now().UnixMilli())
    stage := req.Stage
    if stage == "" { stage = "PROSPECTING" }
    currency := req.Currency
    if currency == "" { currency = "USD" }

    _, err := h.db.Exec(
        `INSERT INTO opportunities (id, account_id, name, value, currency, stage)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        id, req.AccountID, req.Name, req.Value, currency, stage,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    writeJSON(w, http.StatusCreated, map[string]interface{}{
        "id":         id,
        "account_id": req.AccountID,
        "name":       req.Name,
        "value":      req.Value,
        "currency":   currency,
        "stage":      stage,
        "created_at": time.Now(),
    })
}

func (h *OpportunityHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req model.UpdateOpportunityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    res, err := h.db.Exec(
        `UPDATE opportunities SET name=$1, value=$2, currency=$3, stage=$4 WHERE id=$5`,
        req.Name, req.Value, req.Currency, req.Stage, id,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Opportunity not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "updated": true})
}

func (h *OpportunityHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res, err := h.db.Exec(`DELETE FROM opportunities WHERE id = $1`, id)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Opportunity not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "deleted": true})
}