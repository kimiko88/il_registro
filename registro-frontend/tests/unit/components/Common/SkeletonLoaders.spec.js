import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SkeletonTable from '@/components/Common/SkeletonTable.vue'
import SkeletonCard from '@/components/Common/SkeletonCard.vue'


describe('SkeletonTable.vue', () => {
  it('renders default 8 rows and 5 columns', () => {
    const wrapper = mount(SkeletonTable, {
      global: {
        stubs: {
          'q-card': { template: '<div><slot /></div>' },
          'q-card-section': { template: '<div><slot /></div>' }
        }
      }
    })

    expect(wrapper.find('.skeleton-table').exists()).toBe(true)
    const headerCells = wrapper.findAll('.skeleton-header-cell')
    expect(headerCells.length).toBe(5)

    const dataRows = wrapper.findAll('.skeleton-row')
    expect(dataRows.length).toBe(8)

    // Each row has 5 cells
    const firstRowCells = dataRows[0].findAll('.skeleton-cell')
    expect(firstRowCells.length).toBe(5)
  })

  it('renders custom rows and columns when props are provided', () => {
    const wrapper = mount(SkeletonTable, {
      props: {
        rows: 3,
        cols: 2
      }
    })

    const headerCells = wrapper.findAll('.skeleton-header-cell')
    expect(headerCells.length).toBe(2)

    const dataRows = wrapper.findAll('.skeleton-row')
    expect(dataRows.length).toBe(3)
  })
})

describe('SkeletonCard.vue', () => {
  it('renders default 4 skeleton lines without footer', () => {
    const wrapper = mount(SkeletonCard, {
      global: {
        stubs: {
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' }
        }
      }
    })

    // Title line + 4 content lines = 5 skeleton lines
    const lines = wrapper.findAll('.skeleton-line')
    expect(lines.length).toBe(5)
    // Footer shouldn't exist by default
    expect(wrapper.find('.row.q-mt-md.q-gutter-sm').exists()).toBe(false)
  })

  it('renders custom line count and shows footer when showFooter is true', () => {
    const wrapper = mount(SkeletonCard, {
      props: {
        lines: 6,
        showFooter: true,
        cardClass: 'custom-card-class'
      },
      global: {
        stubs: {
          'q-card': { template: '<div class="q-card"><slot /></div>' },
          'q-card-section': { template: '<div class="q-card-section"><slot /></div>' }
        }
      }
    })

    // Title (1) + 6 lines + 2 footer lines = 9 lines
    const lines = wrapper.findAll('.skeleton-line')
    expect(lines.length).toBe(9)
    expect(wrapper.find('.row.q-mt-md.q-gutter-sm').exists()).toBe(true)
  })
})
