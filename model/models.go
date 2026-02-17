package model

import "time"

type Account struct {
    ID            string    `json:"id" db:"id"`
    Name          string    `json:"name" db:"name"`
    Email         string    `json:"email" db:"email"`
    Phone         string    `json:"phone" db:"phone"`
    Industry      string    `json:"industry" db:"industry"`
    SAPCustomerID string    `json:"sap_customer_id" db:"sap_customer_id"`
    Status        string    `json:"status" db:"status"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Opportunity struct {
    ID         string    `json:"id" db:"id"`
    AccountID  string    `json:"account_id" db:"account_id"`
    Name       string    `json:"name" db:"name"`
    Value      float64   `json:"value" db:"value"`
    Currency   string    `json:"currency" db:"currency"`
    Stage      string    `json:"stage" db:"stage"`
    SAPOrderID string    `json:"sap_order_id" db:"sap_order_id"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type Lead struct {
    ID        string    `json:"id" db:"id"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Email     string    `json:"email" db:"email"`
    Company   string    `json:"company" db:"company"`
    Status    string    `json:"status" db:"status"`
    Source    string    `json:"source" db:"source"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Contact struct {
    ID        string    `json:"id" db:"id"`
    AccountID string    `json:"account_id" db:"account_id"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Email     string    `json:"email" db:"email"`
    Phone     string    `json:"phone" db:"phone"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateOpportunityRequest struct {
    AccountID string  `json:"account_id"`
    Name      string  `json:"name"`
    Value     float64 `json:"value"`
    Currency  string  `json:"currency"`
    Stage     string  `json:"stage"`
}

type UpdateLeadRequest struct {
    Status string `json:"status"`
    Source string `json:"source"`
}

type CreateAccountRequest struct {
    ID            string `json:"id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    Phone         string `json:"phone"`
    Industry      string `json:"industry"`
    SAPCustomerID string `json:"sap_customer_id"`
}

type UpdateAccountRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    Industry string `json:"industry"`
    Status   string `json:"status"`
}

type CreateLeadRequest struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Company   string `json:"company"`
    Source    string `json:"source"`
}

type UpdateOpportunityRequest struct {
    Name     string  `json:"name"`
    Value    float64 `json:"value"`
    Currency string  `json:"currency"`
    Stage    string  `json:"stage"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}