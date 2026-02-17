package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "crm-mock/model"
)

type AccountHandler struct { db *sql.DB }

func NewAccountHandler(db *sql.DB) *AccountHandler {
    return &AccountHandler{db: db}
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var a model.Account
    err := h.db.QueryRow(
        `SELECT id, name, email, phone, industry, sap_customer_id, status, created_at
         FROM accounts WHERE id = $1`, id,
    ).Scan(&a.ID, &a.Name, &a.Email, &a.Phone, &a.Industry,
           &a.SAPCustomerID, &a.Status, &a.CreatedAt)

    if err == sql.ErrNoRows {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Account not found", Code: 404})
        return
    }
    writeJSON(w, http.StatusOK, a)
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req model.CreateAccountRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    _, err := h.db.Exec(
        `INSERT INTO accounts (id, name, email, phone, industry, sap_customer_id)
         VALUES ($1, $2, $3, $4, $5, $6)`,
        req.ID, req.Name, req.Email, req.Phone, req.Industry, req.SAPCustomerID,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    writeJSON(w, http.StatusCreated, req)
}

func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var req model.UpdateAccountRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest,
            model.ErrorResponse{Error: "invalid_body", Message: err.Error(), Code: 400})
        return
    }

    res, err := h.db.Exec(
        `UPDATE accounts SET name=$1, email=$2, phone=$3, industry=$4, status=$5 WHERE id=$6`,
        req.Name, req.Email, req.Phone, req.Industry, req.Status, id,
    )
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Account not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "updated": true})
}

func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res, err := h.db.Exec(`DELETE FROM accounts WHERE id = $1`, id)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError,
            model.ErrorResponse{Error: "db_error", Message: err.Error(), Code: 500})
        return
    }

    rows, _ := res.RowsAffected()
    if rows == 0 {
        writeJSON(w, http.StatusNotFound,
            model.ErrorResponse{Error: "not_found", Message: "Account not found", Code: 404})
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "deleted": true})
}
