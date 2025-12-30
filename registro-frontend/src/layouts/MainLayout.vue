<template>
  <q-layout view="lHh Lpr lFf" class="bg-slate-50">
    <q-header class="glass-effect text-slate-900 q-py-xs" :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'">
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          color="primary"
          @click="toggleLeftDrawer"
        />

        <q-toolbar-title class="text-weight-bold text-primary">
          Registro Elettronico
        </q-toolbar-title>

        <div class="text-caption text-grey-7">v0.0.1</div>
        <q-btn flat round dense icon="account_circle" color="primary" class="q-ml-sm" />
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      :class="$q.dark.isActive ? 'bg-dark' : 'bg-white'"
      :width="260"
    >
      <div class="q-pa-md">
        <div class="text-overline text-grey-6 q-mb-sm">MENU</div>
        <q-list padding class="rounded-borders">
          <q-item 
            clickable 
            v-ripple
            tag="a" 
            to="/" 
            active-class="bg-primary text-white rounded-lg shadow-soft"
          >
            <q-item-section avatar>
              <q-icon name="dashboard" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium">Dashboard</q-item-label>
            </q-item-section>
          </q-item>
          
          <!-- Add dynamic menu items here based on role -->

        </q-list>
      </div>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'

const leftDrawerOpen = ref(false)

function toggleLeftDrawer () {
  leftDrawerOpen.value = !leftDrawerOpen.value
}
</script>
