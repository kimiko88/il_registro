import { defineStore } from 'pinia';
import documentService from '../services/documentService';

export const useDocumentsStore = defineStore('documents', {
    state: () => ({
        inbox: [],
        currentDocument: null,
        loading: false,
        error: null,
        pagination: {
            page: 1,
            rowsPerPage: 10,
            rowsNumber: 0
        }
    }),
    actions: {
        async fetchInbox(params = {}) {
            this.loading = true;
            try {
                const response = await documentService.getInbox({
                    page: params.page || this.pagination.page,
                    limit: params.rowsPerPage || this.pagination.rowsPerPage,
                    ...params
                });
                this.inbox = response.data.items;
                this.pagination.rowsNumber = response.data.total;
                this.pagination.page = params.page || this.pagination.page;
                this.pagination.rowsPerPage = params.rowsPerPage || this.pagination.rowsPerPage;
            } catch (err) {
                this.error = err.message;
            } finally {
                this.loading = false;
            }
        },
        async fetchDocument(id) {
            this.loading = true;
            try {
                const response = await documentService.getDocument(id);
                this.currentDocument = response.data;
            } finally {
                this.loading = false;
            }
        },
        async reviewDocument(id, decision, notes) {
            await documentService.reviewDocument(id, decision, notes);
            await this.fetchInbox(); // Refresh
        },
        // Teacher Actions
        async fetchMyDocuments() {
            this.loading = true;
            try {
                // Mock API
                await new Promise(resolve => setTimeout(resolve, 500));
                this.inbox = [ // Reusing inbox state for list, or separate? Let's use 'inbox' as generic list for now or add 'myDocuments' state
                    { id: 'doc1', title: 'PDP - Mario Rossi', status: 'Draft', type: 'PDP', date: '2025-01-20' },
                    { id: 'doc2', title: 'PFI - Sofia Bianchi', status: 'Approved', type: 'PFI', date: '2025-01-10' }
                ];
            } finally {
                this.loading = false;
            }
        },
        async fetchTemplates() {
            // Mock Templates
            return [
                { id: 't1', name: 'PDP Standard', type: 'PDP', content: '<h1>PDP Template</h1>...' },
                { id: 't2', name: 'PFI 2024', type: 'PFI', content: '<h1>PFI Template</h1>...' }
            ];
        },
        async createDocument(docData) {
            // Mock Create
            await new Promise(resolve => setTimeout(resolve, 500));
            this.inbox.unshift({ ...docData, id: 'new', status: 'Draft', date: new Date().toISOString() });
        },

        // Student Actions
        async fetchMyFiles() {
            this.loading = true;
            try {
                // Mock Data
                await new Promise(resolve => setTimeout(resolve, 500));
                this.inbox = [ // Sharing inbox or separate list
                    { id: 1, title: 'Report Card Semester 1', date: '2025-01-15', type: 'Report', extension: 'pdf' },
                    { id: 2, title: 'PDP Signed', date: '2024-11-20', type: 'PDP', extension: 'pdf' }
                ];
            } finally {
                this.loading = false;
            }
        }
    }
});
