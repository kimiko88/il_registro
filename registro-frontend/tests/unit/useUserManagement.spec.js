import { setActivePinia, createPinia } from 'pinia';
import { useUserManagement } from 'src/composables/useUserManagement';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import adminService from 'src/services/adminService';

// Mock Quasar
const mockNotify = vi.fn();
vi.mock('quasar', () => ({
    useQuasar: () => ({
        notify: mockNotify
    })
}));

vi.mock('src/services/adminService');

describe('useUserManagement', () => {
    beforeEach(() => {
        setActivePinia(createPinia());
        mockNotify.mockClear();
    });

    it('fetches users successfully', async () => {
        adminService.getAdmins.mockResolvedValue({
            data: { items: [{ id: 1, name: 'User' }] }
        });

        const { users, fetchUsers, loading } = useUserManagement();

        expect(loading.value).toBe(false);
        await fetchUsers();

        expect(loading.value).toBe(false);
        expect(users.value).toHaveLength(1);
        expect(users.value[0].name).toBe('User');
    });

    it('handles import users mock', async () => {
        const { importUsers } = useUserManagement();

        await importUsers(new File([''], 'test.csv'));

        expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({
            type: 'positive',
            message: expect.any(String)
        }));
    });
});
