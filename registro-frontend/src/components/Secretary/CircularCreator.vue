<template>
  <q-card style="width: min(600px, 95vw); max-width: 95vw;">
    <q-card-section class="row items-center justify-between">
        <div class="text-h6">{{ t('communicationsPage.newCircular') }}</div>
        <q-btn icon="close" flat round dense v-close-popup :aria-label="t('common.close') || 'Chiudi'" />
    </q-card-section>

    <q-card-section>
        <q-form @submit="sendCircular" class="q-gutter-md">
            <q-input v-model="form.title" :label="t('communicationsPage.titleLabel')" outlined dense :rules="[val => !!val || t('common.error')]" />
            
            <div class="text-subtitle2">{{ t('communicationsPage.recipientRole') }}</div>
            <div class="row q-gutter-sm">
                <q-checkbox v-model="form.recipients.teachers" :label="t('usersPage.roleTeachers') || 'Docenti'" />
                <q-checkbox v-model="form.recipients.parents" :label="t('usersPage.roleParents') || 'Genitori'" />
                <q-checkbox v-model="form.recipients.students" :label="t('usersPage.roleStudents') || 'Studenti'" />
                <q-checkbox v-model="form.recipients.staff" :label="t('usersPage.roleStaff') || 'Personale ATA'" />
            </div>

            <q-select
                v-if="form.recipients.students || form.recipients.parents"
                v-model="form.specificClasses"
                multiple
                use-chips
                emit-value
                map-options
                :options="classOptions"
                :label="t('udaPage.classLabel')"
                outlined
                dense
            />

            <div class="q-my-sm">
                <div class="text-subtitle2 q-mb-xs">{{ t('communicationsPage.bodyLabel') }}</div>
                <q-editor v-model="form.content" min-height="150px" />
            </div>

            <q-file v-model="form.attachments" multiple :label="t('communicationsPage.hasAttachment')" outlined dense use-chips>
                <template v-slot:prepend><q-icon name="attach_file" /></template>
            </q-file>

            <div class="row justify-end q-mt-md">
                <q-btn :label="t('common.cancel')" flat v-close-popup color="grey" />
                <q-btn :label="t('communicationsPage.publish')" type="submit" color="primary" class="q-ml-sm" icon="send" :loading="sending" />
            </div>
        </q-form>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useQuasar } from 'quasar'
import { useCommunicationsStore } from '@/stores/communications'
import { useClassesStore } from '@/stores/classes'

const $q = useQuasar()
const { t } = useI18n()
const commStore = useCommunicationsStore()
const classesStore = useClassesStore()
const emit = defineEmits(['sent', 'cancel'])

const sending = ref(false)

const classOptions = computed(() =>
    classesStore.classes.map(c => ({
        label: c.name || `${c.year || ''}${c.section || ''}`.trim() || c.id,
        value: c.id
    }))
)

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
        $q.notify({ type: 'warning', message: t('communicationsPage.recipientsLabel') })
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
        $q.notify({ type: 'positive', message: t('common.success') })
        emit('sent')
    } catch (err) {
        $q.notify({ type: 'negative', message: t('common.error') })
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
