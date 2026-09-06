import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import GradeMatrixGrid from '@/components/Teacher/GradeMatrixGrid.vue';
import api from '@/services/api';

const { mockNotify, mockNotifyStatic } = vi.hoisted(() => ({
  mockNotify: vi.fn(),
  mockNotifyStatic: vi.fn()
}));

vi.mock('quasar', async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    Notify: {
      create: mockNotifyStatic
    },
    useQuasar: () => ({
      notify: mockNotify,
      dark: { isActive: false }
    })
  };
});

vi.mock('@/services/api', () => ({
  default: {
    post: vi.fn().mockResolvedValue({ data: { success: true } })
  }
}));

describe('GradeMatrixGrid.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  const studentsSample = [
    { id: 'std-1', first_name: 'Mario', last_name: 'Rossi' },
    { id: 'std-2', first_name: 'Luigi', last_name: 'Verdi' }
  ];

  const commonStubs = {
    'q-card': { template: '<div><slot /></div>' },
    'q-card-section': { template: '<div><slot /></div>' },
    'q-separator': { template: '<hr />' },
    'q-btn': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
    'q-icon': { template: '<i></i>' },
    'q-tooltip': { template: '<span><slot /></span>' },
    'q-input': {
      props: ['modelValue', 'rules', 'inputClass'],
      template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', Number($event.target.value))" />'
    },
    'CompensativeMeasuresSelector': { template: '<div class="compensative-stub"></div>' }
  };

  it('renders student list correctly with row numbers and names', () => {
    const wrapper = mount(GradeMatrixGrid, {
      props: {
        studentsList: studentsSample,
        subjectId: 'sub-1',
        classId: 'cls-1'
      },
      global: {
        stubs: commonStubs
      }
    });

    expect(wrapper.text()).toContain('Rossi Mario');
    expect(wrapper.text()).toContain('Verdi Luigi');
    expect(wrapper.vm.students.length).toBe(2);
  });

  it('notifies warning if saving with no grades entered', async () => {
    const wrapper = mount(GradeMatrixGrid, {
      props: {
        studentsList: studentsSample,
        subjectId: 'sub-1',
        classId: 'cls-1'
      },
      global: {
        stubs: commonStubs
      }
    });

    await wrapper.vm.saveAllGrades();
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'warning' }));
    expect(api.post).not.toHaveBeenCalled();
  });

  it('blocks saving and shows warning if a grade is out of range (< 1 or > 10)', async () => {
    const wrapper = mount(GradeMatrixGrid, {
      props: {
        studentsList: studentsSample,
        subjectId: 'sub-1',
        classId: 'cls-1'
      },
      global: {
        stubs: commonStubs
      }
    });

    // Enter an invalid grade: 12
    wrapper.vm.students[0].grade_value = 12;

    await wrapper.vm.saveAllGrades();
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'warning' }));
    expect(api.post).not.toHaveBeenCalled();
  });

  it('submits grades successfully when valid and emits saved', async () => {
    const wrapper = mount(GradeMatrixGrid, {
      props: {
        studentsList: studentsSample,
        subjectId: 'sub-1',
        classId: 'cls-1'
      },
      global: {
        stubs: commonStubs
      }
    });

    // Enter valid grades
    wrapper.vm.students[0].grade_value = 8.5;
    wrapper.vm.students[0].notes = 'Ottima prova';
    wrapper.vm.students[1].grade_value = 6;

    await wrapper.vm.saveAllGrades();

    expect(api.post).toHaveBeenCalledWith(
      '/grades/bulk',
      expect.objectContaining({
        subject_id: 'sub-1',
        class_id: 'cls-1',
        grades: expect.arrayContaining([
          expect.objectContaining({ student_id: 'std-1', grade_value: 8.5, description: 'Ottima prova' }),
          expect.objectContaining({ student_id: 'std-2', grade_value: 6 })
        ])
      })
    );
    expect(mockNotify).toHaveBeenCalledWith(expect.objectContaining({ type: 'positive' }));
    expect(wrapper.emitted('saved')).toBeTruthy();
  });

  it('handles keyboard navigation helpers (focusNext / focusPrev)', () => {
    const wrapper = mount(GradeMatrixGrid, {
      props: {
        studentsList: studentsSample,
        subjectId: 'sub-1',
        classId: 'cls-1'
      },
      global: {
        stubs: commonStubs
      }
    });

    const mockFocus1 = vi.fn();
    const mockFocus2 = vi.fn();
    wrapper.vm.inputRefs = [{ focus: mockFocus1 }, { focus: mockFocus2 }];

    wrapper.vm.focusNext(0);
    expect(mockFocus2).toHaveBeenCalled();

    wrapper.vm.focusPrev(1);
    expect(mockFocus1).toHaveBeenCalled();
  });

  it('enqueues to offline queue when offline and emits saved', async () => {
    const originalOnLine = navigator.onLine;
    try {
      Object.defineProperty(navigator, 'onLine', { value: false, configurable: true });
      const wrapper = mount(GradeMatrixGrid, {
        props: {
          studentsList: studentsSample,
          subjectId: 'sub-1',
          classId: 'cls-1'
        },
        global: {
          stubs: commonStubs
        }
      });

      wrapper.vm.students[0].grade_value = 7.5;
      await wrapper.vm.saveAllGrades();

      expect(mockNotifyStatic).toHaveBeenCalledWith(expect.objectContaining({ type: 'warning' }));
      expect(wrapper.emitted('saved')).toBeTruthy();
    } finally {
      Object.defineProperty(navigator, 'onLine', { value: originalOnLine, configurable: true });
    }
  });
});
