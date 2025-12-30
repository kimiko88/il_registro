import { setActivePinia, createPinia } from 'pinia';
import { useAdminStore } from 'src/stores/admin';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import adminService from 'src/services/adminService';

vi.mock('src/services/adminService');

describe('Admin Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
    });

    it('fetches admins', async () => {
        const store = useAdminStore();
        const mockAdmins = [{ id: 1, email: 'admin@test.com' }];

        adminService.getAdmins.mockResolvedValue({
            data: { items: mockAdmins, total: 1 }
        });

        await store.fetchAdmins();

        expect(store.admins).toHaveLength(1);
        expect(store.admins[0].email).toBe('admin@test.com');
        expect(store.loading).toBe(false);
    });

    it('handles error during fetch', async () => {
        const store = useAdminStore();
        adminService.getAdmins.mockRejectedValue(new Error('API Error'));

        await store.fetchAdmins();

        expect(store.error).toBe('API Error');
        expect(store.loading).toBe(false);
    });
});
