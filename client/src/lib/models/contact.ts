export interface Contact {
    id: number;
    name: string;
    linkedInUrl: string;
    credlyInUrl: string;
    dateCreated: Date;
    dateUpdated: Date;
    tags: string[];
    contact: string;
}

export const emptyContact = (): Contact => ({
    id: 0,
    name: '',
    linkedInUrl: '',
    credlyInUrl: '',
    dateCreated: new Date(),
    dateUpdated: new Date(),
    tags: [],
    contact: ''
});