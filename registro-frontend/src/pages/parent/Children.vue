<template>
  <q-page class="q-pa-md bg-slate-50">
    <div class="text-h5 text-weight-bold text-slate-800 q-mb-md">{{ t('nav.children') || 'I Miei Figli' }}</div>

    <div v-if="loading" class="row justify-center">
      <q-spinner color="primary" size="3em" />
    </div>

    <div v-else class="row q-col-gutter-lg">
      <div class="col-12 col-md-6" v-for="child in children" :key="child.id">
        <q-card class="shadow-sm rounded-lg overflow-hidden">
          <q-img src="https://cdn.quasar.dev/img/material.png" height="150px">
            <div class="absolute-bottom text-subtitle1">
              {{ child.first_name }} {{ child.last_name }}
            </div>
          </q-img>

          <q-card-section>
            <div class="row no-wrap items-center">
              <div class="col text-h6 ellipsis">{{ child.school_name }}</div>
            </div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <div class="text-subtitle2 text-grey">{{ t('common.filter') + ' Classe ' + child.class }}</div>
          </q-card-section>

          <q-separator />

          <q-card-actions align="right">
            <q-btn flat round icon="email" color="primary">
              <q-tooltip>{{ t('common.search') }}</q-tooltip>
            </q-btn>
            <q-btn flat :label="t('nav.myGrades')" color="primary" to="/parent/grades" @click="selectChild(child.id)" />
          </q-card-actions>
          
          <q-expansion-item icon="info" :label="t('documentsPage.title')" header-class="text-primary">
            <q-card>
              <q-card-section>
                <p><strong>Email:</strong> segreteria@scuola.it</p>
              </q-card-section>
            </q-card>
          </q-expansion-item>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useParentStore } from '@/stores/parent'
import { storeToRefs } from 'pinia'

const { t } = useI18n()
const parentStore = useParentStore()
const { children, loading } = storeToRefs(parentStore)
const { fetchChildren, selectChild } = parentStore

onMounted(() => {
  if (children.value.length === 0) {
    fetchChildren()
  }
})
</script>
