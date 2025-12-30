<template>
  <q-page class="q-pa-md" style="height: calc(100vh - 50px)">
    <div class="row full-height q-col-gutter-md">
        <!-- Thread List -->
        <div class="col-12 col-md-4 full-height">
            <q-card class="full-height column">
                <q-card-section class="bg-primary text-white">
                    <div class="text-h6">Inbox <q-badge color="orange" v-if="store.unreadCount">{{ store.unreadCount }}</q-badge></div>
                </q-card-section>
                <q-scroll-area class="col">
                    <q-list separator>
                        <q-item 
                            v-for="thread in store.threads" 
                            :key="thread.id" 
                            clickable 
                            @click="selectThread(thread.id)"
                            :active="store.activeThreadId === thread.id"
                            active-class="bg-blue-1"
                        >
                            <q-item-section avatar>
                                <q-avatar color="grey" text-color="white" icon="person" />
                            </q-item-section>
                            <q-item-section>
                                <q-item-label :class="{'text-weight-bold': thread.unread}">{{ thread.subject }}</q-item-label>
                                <q-item-label caption>{{ thread.recipient }}</q-item-label>
                                <q-item-label caption lines="1">{{ thread.lastMessage }}</q-item-label>
                            </q-item-section>
                            <q-item-section side top>
                                <q-item-label caption>{{ thread.date }}</q-item-label>
                                <q-icon name="circle" color="blue" size="xs" v-if="thread.unread" />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-scroll-area>
            </q-card>
        </div>

        <!-- Active Thread -->
        <div class="col-12 col-md-8 full-height">
            <q-card class="full-height column" v-if="store.activeThread">
                <q-card-section class="bg-grey-1 border-bottom">
                    <div class="text-subtitle1">{{ store.activeThread.subject }}</div>
                    <div class="text-caption">To: {{ store.activeThread.recipient }}</div>
                </q-card-section>
                <div class="col full-height" style="overflow: hidden">
                    <MessageThread :thread-id="store.activeThread.id" :messages="store.activeThread.messages" />
                </div>
            </q-card>
            <div v-else class="full-height flex flex-center text-grey">
                Select a conversation
            </div>
        </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useCommunicationsStore } from 'src/stores/communications';
import MessageThread from 'src/components/Teacher/MessageThread.vue';

const store = useCommunicationsStore();

const selectThread = (id) => {
    store.fetchMessages(id);
};

onMounted(() => {
    store.fetchThreads();
});
</script>
