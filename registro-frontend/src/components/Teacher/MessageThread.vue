<template>
  <div class="column full-height">
    <!-- Messages Area -->
    <q-scroll-area class="col q-pa-md">
        <div v-for="msg in messages" :key="msg.id" :class="['q-mb-md row', msg.from === 'Me' ? 'justify-end' : 'justify-start']">
            <q-chat-message
                :name="msg.from"
                :text="[msg.text]"
                :sent="msg.from === 'Me'"
                :stamp="msg.date"
                :bg-color="msg.from === 'Me' ? 'primary' : 'grey-3'"
                :text-color="msg.from === 'Me' ? 'white' : 'black'"
            />
        </div>
    </q-scroll-area>

    <!-- Input Area -->
    <div class="q-pa-md bg-grey-2">
        <div class="row q-gutter-sm">
            <q-input v-model="text" outlined dense class="col" placeholder="Type a message..." @keyup.enter="send" :disable="loading" />
            <q-btn round flat icon="send" color="primary" @click="send" :loading="loading" />
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useMessages } from 'src/composables/useMessages';

const props = defineProps(['threadId', 'messages']);
const { sendMessage, sending: loading } = useMessages();
const text = ref('');

const send = async () => {
    if (await sendMessage(props.threadId, text.value)) {
        text.value = '';
    }
};
</script>
