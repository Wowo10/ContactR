export interface User {
    id: string;
    name: string;
    email: string;
    valid_until: Date;
    is_valid: boolean;
    is_admin: boolean;
}

export interface LoginData {
    user_id: string;
    user_email: string;
    is_valid: boolean;
    is_admin: boolean;
}