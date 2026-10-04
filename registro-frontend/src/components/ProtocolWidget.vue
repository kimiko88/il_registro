<template>
  <div class="inline-block">
    <q-btn
      :color="isProtocolled ? 'teal-8' : 'primary'"
      :icon="isProtocolled ? 'verified' : 'mark_email_read'"
      :label="isProtocolled ? `Prot. N. ${protocolNumber}` : t('agidProtocol.widgetTitle')"
      size="sm"
      :flat="isProtocolled"
      :unelevated="!isProtocolled"
      no-caps
      @click="onClick"
    />

    <q-dialog v-model="showDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ t('agidProtocol.widgetTitle') }}</div>
        </q-card-section>
        <q-card-section class="q-pt-md">
          <p class="text-caption">{{ t('agidProtocol.widgetHelp') }}</p>
          <q-input v-model="form.subject" :label="t('agidProtocol.subject') + ' *'" outlined dense class="q-mb-sm" />
          <q-input v-model="form.sender" :label="t('agidProtocol.sender') + ' *'" outlined dense class="q-mb-sm" />
          <q-input v-model="form.recipient" :label="t('agidProtocol.recipient') + ' *'" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="t('common.cancel')" v-close-popup />
          <q-btn color="primary" :label="t('agidProtocol.newEntry')" unelevated :loading="submitting" @click="submitProtocol" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { protocolService } from '@/services/protocolService'

export default defineComponent({
  name: 'ProtocolWidget',
  props: {
    entityType: { type: String, required: true },
    entityId: { type: String, required: true },
    defaultSubject: { type: String, default: '' },
    defaultSender: { type: String, default: 'Istituto Scolastico' },
    defaultRecipient: { type: String, default: 'Famiglie / Personale' }
  },
  emits: ['protocolled'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const showDialog = ref(false)
    const submitting = ref(false)
    const isProtocolled = ref(false)
    const protocolNumber = ref('')

    const form = ref({
      subject: props.defaultSubject,
      sender: props.defaultSender,
      recipient: props.defaultRecipient
    })

    const onClick = () => {
      if (!isProtocolled.value) {
        form.value.subject = props.defaultSubject
        showDialog.value = true
      }
    }

    const submitProtocol = async () => {
      submitting.value = true
      try {
        const res = await protocolService.protocolDocument({
          subject: form.value.subject,
          sender: form.value.sender,
          recipient: form.value.recipient,
          entity_type: props.entityType,
          entity_id: props.entityId
        })
        isProtocolled.value = true
        protocolNumber.value = String(res.protocol.protocol_number).padStart(7, '0')
        showDialog.value = false
        emit('protocolled', res.protocol)
      } catch (e) {
        console.error('Errore protocollo rapido:', e)
      } finally {
        submitting.value = false
      }
    }

    return {
      t,
      showDialog,
      submitting,
      isProtocolled,
      protocolNumber,
      form,
      onClick,
      submitProtocol
    }
  }
})
</script>
