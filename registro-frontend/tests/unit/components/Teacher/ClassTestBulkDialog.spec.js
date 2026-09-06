import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import ClassTestBulkDialog from '@/components/Teacher/ClassTestBulkDialog.vue';
import { gradeService } from '@/services/gradeService';

const { mockNotify } = vi.hoisted(() => ({
  mockNotify: vi.fn()
}));

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    useQuasar: () => ({
      notify: mockNotify,
      dark: { isActive: false }
    })
  };
});

vi.mock('@/services/gradeService', () => ({
  gradeService: {
    createTestWithGrades: vi.fn().mockResolvedValue({ data: { id: 'test-1' } })
  }
}));

describe('ClassTestBulkDialog.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const studentsSample = [
    { student_id: 's1', full_name: 'Mario Rossi', grade_value: null, notes: '' },
    { student_id: 's2', full_name: 'Luigi Verdi', grade_value: null, notes: '' }
  ];

  const commonStubs = {
    'q-dialog': { template: '<div><slot /></div>' },
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-card-actions': { template: '<div><slot /></div>' },
    'q-form': {
      template: '<form @submit.prevent="$emit(\'submit\')"><slot /></form>',
      methods: { validate: () => Promise.resolve(true) }
    },
    'q-input': { template: '<div class="q-input-stub"><slot /></div>' },
    'q-select': { template: '<div class="q-select-stub"><slot /></div>' },
    'q-scroll-area': { template: '<div><slot /></div>' },
    'q-list': { template: '<div><slot /></div>' },
    'q-item': { template: '<div><slot /></div>' },
    'q-item-section': { template: '<div><slot /></div>' },
    'q-item-label': { template: '<div><slot /></div>' },
    'q-linear-progress': { template: '<div></div>' },
    'q-banner': { template: '<div><slot /></div>' },
    'q-btn': { template: '<button><slot /></button>' },
    'q-space': { template: '<div></div>' },
    'q-icon': { template: '<i></i>' }
  };

  it('renders create dialog header when isEdit is false', () => {
    const wrapper = mount(ClassTestBulkDialog, {
      props: {
        modelValue: true,
        isEdit: false,
        classId: 'c1',
        subjectId: 'sub1',
        testData: {
          title: 'Verifica Geometria',
          date: '2026-09-10',
          evaluationType: 'Scritto',
          grades: studentsSample
        }
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    expect(wrapper.text()).toContain('Crea Nuova Verifica');
  });

  it('renders edit dialog header when isEdit is true', () => {
    const wrapper = mount(ClassTestBulkDialog, {
      props: {
        modelValue: true,
        isEdit: true,
        classId: 'c1',
        subjectId: 'sub1',
        testData: {
          id: 'test-100',
          title: 'Verifica Algebra',
          date: '2026-09-12',
          evaluationType: 'Scritto',
          grades: studentsSample
        }
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    expect(wrapper.text()).toContain('Modifica Verifica in Blocco');
  });

  it('marks all students as absent when markAllAbsent is called', async () => {
    const wrapper = mount(ClassTestBulkDialog, {
      props: {
        modelValue: true,
        isEdit: false,
        classId: 'c1',
        subjectId: 'sub1',
        testData: {
          title: 'Verifica',
          date: '2026-09-10',
          evaluationType: 'Scritto',
          grades: [
            { student_id: 's1', full_name: 'Mario Rossi', grade_value: 8, notes: '' },
            { student_id: 's2', full_name: 'Luigi Verdi', grade_value: 7, notes: '' }
          ]
        }
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    wrapper.vm.markAllAbsent();
    expect(wrapper.vm.localForm.grades[0].grade_value).toBeNull();
    expect(wrapper.vm.localForm.grades[0].notes).toBe('Assente');
    expect(wrapper.vm.localForm.grades[1].grade_value).toBeNull();
    expect(wrapper.vm.localForm.grades[1].notes).toBe('Assente');
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'info' }));
  });

  it('submits create test successfully and emits saved', async () => {
    const wrapper = mount(ClassTestBulkDialog, {
      props: {
        modelValue: true,
        isEdit: false,
        classId: 'c1',
        subjectId: 'sub1',
        testData: {
          title: 'Verifica Storia',
          date: '2026-09-10',
          evaluationType: 'Scritto',
          grades: [
            { student_id: 's1', full_name: 'Mario Rossi', grade_value: 8, notes: 'Bravo' }
          ]
        }
      },
      global: {
        plugins: [createTestingPinia({ stubActions: false })],
        stubs: commonStubs
      }
    });

    await wrapper.vm.handleSubmit();
    expect(gradeService.createTestWithGrades).toHaveBeenCalledWith(
      expect.objectContaining({
        class_id: 'c1',
        subject_id: 'sub1',
        title: 'Verifica Storia',
        grades: expect.arrayContaining([
          expect.objectContaining({ student_id: 's1', grade_value: 8, notes: 'Bravo' })
        ])
      })
    );
    expect(wrapper.emitted('update:modelValue')).toContainEqual([false]);
    expect(wrapper.emitted('saved')).toBeTruthy();
  });
});
