<template>
  <q-dialog v-model="show" full-width full-height>
    <q-card class="column no-wrap rounded-xl overflow-hidden glass-card">
      <q-card-section class="bg-gradient-premium text-white row items-center q-pa-lg">
        <div class="text-h5 text-weight-bold text-outfit">Gestione Template Documenti</div>
        <q-space />
        <q-btn icon="add" label="Nuovo Template" color="white" text-color="primary" class="rounded-lg q-mr-sm" @click="openCreate" />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="col q-pa-lg scroll">
        <div v-if="loading" class="row justify-center q-pa-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>
        
        <div v-else-if="templates.length === 0" class="column items-center justify-center q-pa-xl text-grey-6">
          <q-icon name="description" size="80px" class="q-mb-md" />
          <div class="text-h6">Nessun template configurato</div>
          <p>Crea un template per velocizzare la generazione di documenti e circolari.</p>
        </div>

        <div v-else class="row q-col-gutter-lg">
          <div v-for="tpl in templates" :key="tpl.id" class="col-12 col-md-4 col-lg-3">
            <q-card class="rounded-xl shadow-soft h-full hover-card border-slate-100">
              <q-card-section>
                <div class="row items-center justify-between q-mb-sm">
                  <q-chip :color="getTypeColor(tpl.type)" text-color="white" dense class="text-weight-bold">
                    {{ tpl.type }}
                  </q-chip>
                  <div class="row q-gutter-xs">
                    <q-btn flat round dense icon="edit" color="primary" size="sm" @click="openEdit(tpl)" />
                    <q-btn flat round dense icon="delete" color="negative" size="sm" @click="confirmDelete(tpl)" />
                  </div>
                </div>
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-xs">{{ tpl.name }}</div>
                <div class="text-caption text-slate-500 ellipsis-3-lines" v-html="tpl.content"></div>
              </q-card-section>
            </q-card>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Create/Edit Dialog -->
    <q-dialog v-model="showEditDialog" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="column no-wrap">
        <q-card-section class="row items-center q-pa-md border-b">
          <div class="text-h6 text-weight-bold">{{ editingTemplate?.id ? 'Modifica Template' : 'Nuovo Template' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="col q-pa-lg scroll">
          <div class="max-w-4xl mx-auto">
            <q-form @submit="saveTemplate" class="q-gutter-md">
              <div class="row q-col-gutter-md">
                <div class="col-8">
                  <q-input v-model="form.name" label="Nome Template" outlined :rules="[val => !!val || 'Obbligatorio']" />
                </div>
                <div class="col-4">
                  <q-select 
                    v-model="form.type" 
                    :options="['circular', 'certificate', 'report', 'other']" 
                    label="Tipo" 
                    outlined 
                  />
                </div>
              </div>

              <div class="text-subtitle2 q-mb-xs">Contenuto Template</div>
              <div class="text-caption text-grey-7 q-mb-sm">
                Puoi usare placeholder come <code>{student_name}</code>, <code>{date}</code>, <code>{school_name}</code>.
              </div>
              
              <q-editor
                v-model="form.content"
                min-height="400px"
                flat
                bordered
                class="rounded-lg overflow-hidden"
                :toolbar="[
                  ['bold', 'italic', 'strike', 'underline'],
                  ['token', 'hr', 'link', 'custom_btn'],
                  ['print', 'fullscreen'],
                  [
                    {
                      label: $q.lang.editor.formatting,
                      icon: $q.iconSet.editor.formatting,
                      list: 'no-icons',
                      options: ['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'code']
                    },
                    {
                      label: $q.lang.editor.fontSize,
                      icon: $q.iconSet.editor.fontSize,
                      fixedLabel: true,
                      fixedIcon: true,
                      list: 'no-icons',
                      options: ['size-1', 'size-2', 'size-3', 'size-4', 'size-5', 'size-6', 'size-7']
                    }
                  ],
                  ['quote', 'unordered', 'ordered', 'outdent', 'indent'],
                  ['undo', 'redo'],
                  ['viewsource']
                ]"
              />

              <div class="row justify-end q-mt-xl q-gutter-sm">
                <q-btn label="Annulla" flat v-close-popup />
                <q-btn :label="editingTemplate?.id ? 'Aggiorna Template' : 'Crea Template'" type="submit" color="primary" class="q-px-xl rounded-lg" :loading="saving" />
              </div>
            </q-form>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useQuasar } from 'quasar'
import api from 'src/services/api'

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(['update:modelValue', 'templates-updated'])

const $q = useQuasar()
const show = ref(false)
const loading = ref(false)
const saving = ref(false)
const templates = ref([])
const showEditDialog = ref(false)
const editingTemplate = ref(null)

const form = reactive({
  name: '',
  type: 'circular',
  content: ''
})

watch(() => props.modelValue, (val) => {
  show.value = val
  if (val) fetchTemplates()
})

watch(show, (val) => {
  emit('update:modelValue', val)
})

const fetchTemplates = async () => {
  loading.value = true
  try {
    const res = await api.get('/documents/template')
    templates.value = res.data || []
  } catch (e) {
    console.error('Error fetching templates:', e)
    $q.notify({ type: 'negative', message: 'Errore caricamento template' })
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editingTemplate.value = null
  form.name = ''
  form.type = 'circular'
  form.content = ''
  showEditDialog.value = true
}

const openEdit = (tpl) => {
  editingTemplate.value = tpl
  form.name = tpl.name
  form.type = tpl.type
  form.content = tpl.content
  showEditDialog.value = true
}

const saveTemplate = async () => {
  saving.value = true
  try {
    if (editingTemplate.value) {
      await api.patch(`/documents/template/${editingTemplate.value.id}`, form)
      $q.notify({ type: 'positive', message: 'Template aggiornato' })
    } else {
      await api.post('/documents/template', form)
      $q.notify({ type: 'positive', message: 'Template creato' })
    }
    showEditDialog.value = false
    fetchTemplates()
    emit('templates-updated')
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante il salvataggio' })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (tpl) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Sei sicuro di voler eliminare il template "${tpl.name}"?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await api.delete(`/documents/template/${tpl.id}`)
      $q.notify({ type: 'positive', message: 'Template eliminato' })
      fetchTemplates()
      emit('templates-updated')
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore durante l\'eliminazione' })
    }
  })
}

const getTypeColor = (type) => {
  switch (type) {
    case 'circular': return 'indigo-6';
    case 'certificate': return 'emerald-6';
    case 'report': return 'amber-7';
    default: return 'slate-5';
  }
}
</script>

<style scoped>
.hover-card {
  transition: all 0.3s ease;
}
.hover-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px -8px rgba(0,0,0,0.1) !important;
  border-color: var(--q-primary);
}
.border-slate-100 {
  border: 1px solid #f1f5f9;
}
.max-w-4xl {
  max-width: 56rem;
}
.mx-auto {
  margin-left: auto;
  margin-right: auto;
}
.border-b {
  border-bottom: 1px solid #edf2f7;
}
.h-full {
  height: 100%;
}
</style>
