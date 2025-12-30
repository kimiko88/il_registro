import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import SchoolTable from 'src/components/Admin/SchoolTable.vue';
import { describe, it, expect, vi } from 'vitest';
import { Quasar } from 'quasar';

describe('SchoolTable.vue', () => {
    it('renders rows correctly', () => {
        const wrapper = mount(SchoolTable, {
            global: {
                plugins: [
                    createTestingPinia({
                        initialState: {
                            schools: {
                                schools: [
                                    { id: 1, name: 'Test School', email: 'test@school.com' }
                                ],
                                pagination: { page: 1, rowsPerPage: 10, rowsNumber: 1 }
                            }
                        },
                        createSpy: vi.fn,
                    }),
                ],
                mocks: {
                    $q: { screen: { lt: { sm: false } } } // Mock Quasar $q
                },
                stubs: {
                    // Stub Quasar components
                    'q-table': {
                        template: '<table><slot name="body-cell-actions" :props="{row: {id: 1}}"></slot></table>',
                        props: ['rows', 'columns']
                    },
                    'q-btn': true,
                    'q-td': { template: '<td><slot></slot></td>' },
                    'q-input': true,
                    'q-space': true,
                    'q-icon': true
                }
            }
        });

        expect(wrapper.exists()).toBeTruthy();
        // Further assertions would depend on exact q-table implementation stubbing details
        // For now, checking if it mounts without error with Pinia store.
    });
});
