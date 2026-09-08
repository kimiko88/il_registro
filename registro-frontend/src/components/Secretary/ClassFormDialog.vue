<template>
  <q-dialog :model-value="modelValue" persistent class="premium-dialog" @update:model-value="$emit('update:modelValue', $event)">
    <q-card style="width: min(480px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden shadow-24">
      <q-card-section class="bg-gradient-primary text-white row items-center q-pa-lg">
        <div class="text-h6 text-weight-bold">{{ isEdit ? t('secretaryClasses.editClass') : t('secretaryClasses.newClass') }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close')" />
      </q-card-section>

      <q-card-section class="q-pa-xl">
        <q-form @submit="saveClass" class="q-gutter-y-lg">
          <div class="row q-col-gutter-lg">
            <div class="col-8">
              <q-input
                v-model="form.name"
                :label="t('secretaryClasses.formName')"
                outlined
                :rules="[val => !!val || t('common.required')]"
              />
            </div>
            <div class="col-4">
              <q-input
                v-model="form.section"
                :label="t('secretaryClasses.formSection')"
                outlined
                :rules="[val => !!val || t('common.required')]"
              />
            </div>
          </div>

          <q-input v-model="form.articolazione" :label="t('secretaryClasses.formArticolazione')" outlined />
          <q-input v-model="form.location" :label="t('secretaryClasses.formLocation')" outlined />

          <q-select
            v-model="form.academic_year"
            :options="academicYearOptions"
            :label="t('secretaryClasses.academicYear')"
            outlined
            :rules="[val => !!val || t('common.required')]"
          />

          <q-select
            v-model="form.coordinator_id"
            :options="teacherUserOptions"
            :label="t('secretaryClasses.coordinator')"
            outlined
            emit-value
            map-options
            clearable
          />

          <div class="row justify-end q-mt-xl q-gutter-sm">
            <q-btn :label="t('common.cancel')" flat v-close-popup color="grey-7" />
            <q-btn
              :label="isEdit ? t('common.update') : t('secretaryClasses.createClass')"
              type="submit"
              color="primary"
              class="q-px-xl rounded-lg shadow-sm"
              :loading="saving"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useClassesStore } from '@/stores/classes'
import { useAuthStore } from '@/stores/auth'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  academicYearOptions: {
    type: Array,
    default: () => []
  },
  teacherUserOptions: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'saved'])

const { t } = useI18n()
const $q = useQuasar()
const classesStore = useClassesStore()
const authStore = useAuthStore()
const saving = ref(false)

const form = reactive({
  id: null,
  name: '',
  section: '',
  articolazione: '',
  location: '',
  academic_year: '',
  coordinator_id: ''
})

// Sync form with initialData when dialog opens
watch(() => props.modelValue, (open) => {
  if (open && props.initialData) {
    Object.assign(form, {
      id: null,
      name: '',
      section: '',
      articolazione: '',
      location: '',
      academic_year: '',
      coordinator_id: '',
      ...props.initialData
    })
  }
})

const saveClass = async () => {
  saving.value = true
  try {
    const payload = { ...form, school_id: authStore.user.school_id }
    if (props.isEdit) {
      await classesStore.updateClass(form.id, payload)
      $q.notify({ type: 'positive', message: t('secretaryClasses.classUpdated') })
    } else {
      await classesStore.createClass(payload)
      $q.notify({ type: 'positive', message: t('secretaryClasses.classCreated') })
    }
    emit('update:modelValue', false)
    emit('saved')
  } catch {
    $q.notify({ type: 'negative', message: t('secretaryClasses.saveError') })
  } finally {
    saving.value = false
  }
}
</script>
