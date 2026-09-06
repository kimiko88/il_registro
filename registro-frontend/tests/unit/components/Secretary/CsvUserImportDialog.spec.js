import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import CsvUserImportDialog from '@/components/Secretary/CsvUserImportDialog.vue';
import { exportFile } from 'quasar';
import { useUsersStore } from '@/stores/users';

const { mockNotify } = vi.hoisted(() => ({
  mockNotify: vi.fn()
}));

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    exportFile: vi.fn().mockReturnValue(true),
    useQuasar: () => ({
      notify: mockNotify,
      dark: { isActive: false }
    })
  };
});

describe('CsvUserImportDialog.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const commonStubs = {
    'q-dialog': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-stepper': { template: '<div><slot /></div>' },
    'q-step': { template: '<div class="q-step"><slot /></div>' },
    'q-select': { template: '<div class="q-select"><slot /></div>' },
    'q-toggle': { template: '<div class="q-toggle"><slot /></div>' },
    'q-file': { template: '<div class="q-file"><slot /></div>' },
    'q-table': { template: '<div class="q-table"><slot /></div>' },
    'q-banner': { template: '<div class="q-banner"><slot /></div>' },
    'q-list': { template: '<div><slot /></div>' },
    'q-item': { template: '<div><slot /></div>' },
    'q-item-section': { template: '<div><slot /></div>' },
    'q-item-label': { template: '<div><slot /></div>' },
    'q-btn': { template: '<button><slot /></button>' },
    'q-icon': { template: '<i></i>' }
  };

  it('renders dialog and downloads template', () => {
    const wrapper = mount(CsvUserImportDialog, {
      props: {
        modelValue: true,
        classOptions: [{ label: '1A', value: 'c1' }]
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    expect(wrapper.text()).toContain('Importazione Utenti da CSV');
    wrapper.vm.downloadCsvTemplate();
    expect(exportFile).toHaveBeenCalledWith(
      expect.stringContaining('template_import_'),
      expect.any(String),
      'text/csv'
    );
  });

  it('advances steps and runs bulk import successfully', async () => {
    const wrapper = mount(CsvUserImportDialog, {
      props: {
        modelValue: true,
        classOptions: [{ label: '1A', value: 'c1' }]
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    const store = useUsersStore();
    store.bulkImportUsers = vi.fn().mockResolvedValue({
      imported: 5,
      skipped: 1,
      errors: [{ row: 2, reason: 'Email non valida' }]
    });

    wrapper.vm.importStep = 2;
    wrapper.vm.importFile = new File(['cognome,nome\nRossi,Mario'], 'users.csv', { type: 'text/csv' });

    await wrapper.vm.runBulkImport();

    expect(store.bulkImportUsers).toHaveBeenCalled();
    expect(wrapper.vm.importStep).toBe(3);
    expect(wrapper.vm.importResult.imported).toBe(5);

    // Download error report
    wrapper.vm.downloadErrorReport();
    expect(exportFile).toHaveBeenCalledWith(
      'report_errori_import.csv',
      expect.stringContaining('Email non valida'),
      'text/csv'
    );

    // Close modal
    wrapper.vm.closeImportModal();
    expect(wrapper.emitted('update:modelValue')).toContainEqual([false]);
    expect(wrapper.emitted('imported')).toBeTruthy();
  });
});
