<template>
  <q-card style="min-width: 600px">
    <q-card-section>
        <div class="text-h6">Crea Nuova Circolare</div>
    </q-card-section>

    <q-card-section>
        <q-form @submit="sendCircular" class="q-gutter-md">
            <q-input v-model="form.title" label="Oggetto" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            
            <div class="text-subtitle2">Destinatari</div>
            <div class="row q-gutter-sm">
                <q-checkbox v-model="form.recipients.teachers" label="Docenti" />
                <q-checkbox v-model="form.recipients.parents" label="Genitori" />
                <q-checkbox v-model="form.recipients.students" label="Studenti" />
                <q-checkbox v-model="form.recipients.staff" label="Personale ATA" />
            </div>

            <q-select
                v-if="form.recipients.students || form.recipients.parents"
                v-model="form.specificClasses"
                multiple
                use-chips
                :options="classOptions"
                label="Limita a classi (opzionale)"
                outlined
                dense
                hint="Lascia vuoto per inviare a tutte le classi"
            />

            <div class="q-my-sm">
                <div class="text-subtitle2 q-mb-xs">Contento</div>
                <q-editor v-model="form.content" min-height="150px" />
            </div>

            <q-file v-model="form.attachments" multiple label="Allegati" outlined dense use-chips>
                <template v-slot:prepend><q-icon name="attach_file" /></template>
            </q-file>

            <div class="row justify-end q-mt-md">
                <q-btn label="Annulla" flat v-close-popup color="grey" />
                <q-btn label="Invia Comunicazione" type="submit" color="primary" class="q-ml-sm" icon="send" :loading="sending" />
            </div>
        </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useCommunicationsStore } from 'src/stores/communications'
import { useClassesStore } from 'src/stores/classes'

const $q = useQuasar()
const commStore = useCommunicationsStore()
const classesStore = useClassesStore()
const emit = defineEmits(['sent', 'cancel'])

const sending = ref(false)

const classOptions = computed(() => classesStore.classes.map(c => `${c.name}${c.section}`))

const form = reactive({
    title: '',
    content: '',
    recipients: {
        teachers: false,
        parents: false,
        students: false,
        staff: false
    },
    specificClasses: [],
    attachments: []
})

onMounted(async () => {
    if (classesStore.classes.length === 0) {
        await classesStore.fetchClasses()
    }
})

const sendCircular = async () => {
    if (!form.recipients.teachers && !form.recipients.parents && !form.recipients.students && !form.recipients.staff) {
        $q.notify({ type: 'warning', message: 'Seleziona almeno un gruppo di destinatari' })
        return
    }

    sending.value = true
    try {
        await commStore.sendMessage({
            title: form.title,
            content: form.content,
            recipients: form.recipients,
            specific_classes: form.specificClasses,
            type: 'circular'
        })
        $q.notify({ type: 'positive', message: 'Circolare inviata correttamente' })
        emit('sent')
    } catch (err) {
        $q.notify({ type: 'negative', message: 'Errore durante l\'invio' })
    } finally {
        sending.value = false
    }
}

defineExpose({
    form,
    sendCircular,
    sending
})
</script>
