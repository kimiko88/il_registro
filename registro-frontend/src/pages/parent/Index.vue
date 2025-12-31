<template>
  <q-page class="q-pa-md bg-grey-1">
    
    <!-- Child Switcher Header -->
    <q-card class="bg-primary text-white q-mb-md">
        <q-card-section>
            <div class="row items-center justify-between">
                <div>
                     <div class="text-caption opacity-80">Stai visualizzando il profilo di:</div>
                     <q-btn-dropdown flat :label="parentStore.selectedChild?.name || 'Seleziona Figlio'" icon="face" size="lg" no-caps>
                         <q-list>
                             <q-item v-for="child in parentStore.children" :key="child.id" clickable v-close-popup @click="parentStore.selectChild(child)">
                                 <q-item-section avatar><q-avatar icon="face" color="primary" text-color="white" /></q-item-section>
                                 <q-item-section>
                                     <q-item-label>{{ child.name }}</q-item-label>
                                     <q-item-label caption>{{ child.class }} - {{ child.school }}</q-item-label>
                                 </q-item-section>
                             </q-item>
                         </q-list>
                     </q-btn-dropdown>
                </div>
                <div>
                    <div class="text-h6">{{ parentStore.selectedChild?.class }}</div>
                </div>
            </div>
        </q-card-section>
    </q-card>

    <div v-if="parentStore.selectedChild">
        <!-- Overview Stats -->
        <div class="row q-col-gutter-sm q-mb-lg">
            <div class="col-6 col-sm-3">
                <q-card class="text-center q-py-sm">
                    <div class="text-h5 text-weight-bold text-primary">7.8</div>
                    <div class="text-caption text-grey">Media Voti</div>
                </q-card>
            </div>
            <div class="col-6 col-sm-3">
                <q-card class="text-center q-py-sm">
                     <div class="text-h5 text-weight-bold text-red">3</div>
                     <div class="text-caption text-grey">Assenze</div>
                </q-card>
            </div>
            <div class="col-6 col-sm-3">
                <q-card class="text-center q-py-sm">
                     <div class="text-h5 text-weight-bold text-orange">2</div>
                     <div class="text-caption text-grey">Note</div>
                </q-card>
            </div>
            <div class="col-6 col-sm-3">
                <q-card class="text-center q-py-sm">
                     <div class="text-h5 text-weight-bold text-blue">1</div>
                     <div class="text-caption text-grey">Avvisi</div>
                </q-card>
            </div>
        </div>
        
        <!-- Navigation Grid -->
        <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
                <!-- Latest Grades -->
                <q-card class="q-mb-md">
                    <q-card-section class="row items-center justify-between">
                        <div class="text-h6">Ultimi Voti</div>
                        <q-btn flat label="Vedi Tutti" color="primary" to="/parent/grades" />
                    </q-card-section>
                    <q-list separator>
                        <q-item>
                            <q-item-section>
                                <q-item-label>Matematica</q-item-label>
                                <q-item-label caption>Compito in classe</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <q-badge color="green" label="7.5" />
                            </q-item-section>
                        </q-item>
                        <q-item>
                            <q-item-section>
                                <q-item-label>Storia</q-item-label>
                                <q-item-label caption>Interrogazione</q-item-label>
                            </q-item-section>
                            <q-item-section side>
                                <q-badge color="orange" label="6+" />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-card>
            </div>

            <div class="col-12 col-md-6">
                 <!-- Actions -->
                 <div class="row q-col-gutter-sm">
                     <div class="col-6"><q-btn push color="white" text-color="black" icon="event" label="Colloqui" class="full-width q-py-lg" to="/parent/colloqui" /></div>
                     <div class="col-6"><q-btn push color="white" text-color="black" icon="check_circle" label="Assenze" class="full-width q-py-lg" to="/parent/attendance" /></div>
                     <div class="col-6"><q-btn push color="white" text-color="black" icon="campaign" label="Avvisi" class="full-width q-py-lg" to="/parent/communications" /></div>
                     <div class="col-6"><q-btn push color="white" text-color="black" icon="euro" label="Pagamenti" class="full-width q-py-lg" to="/parent/payments" /></div>
                 </div>
            </div>
        </div>
    </div>

  </q-page>
</template>

<script setup>
import { onMounted } from 'vue';
import { useParentStore } from 'src/stores/parent';

const parentStore = useParentStore();

onMounted(() => {
    parentStore.fetchChildren();
});
</script>

<style scoped>
.opacity-80 { opacity: 0.8; }
</style>
