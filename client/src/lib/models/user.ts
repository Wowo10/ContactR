export interface User {
    user_email: string;
    user_id: string;
    id: string; //TODO: unify this
    name: string;
    email: string;
    valid_until: Date;
    is_valid: boolean;
    is_admin: boolean;
}