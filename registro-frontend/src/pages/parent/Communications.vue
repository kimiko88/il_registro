<template>
  <q-page class="q-pa-md">
    <h1 class="text-h4 q-my-none q-mb-md">Communications</h1>
    
    <q-list bordered separator>
        <q-item v-for="n in 3" :key="n" clickable>
            <q-item-section avatar>
                <q-icon name="campaign" color="primary" />
            </q-item-section>
            <q-item-section>
                <q-item-label>School Closure - Holiday</q-item-label>
                <q-item-label caption>The school will be closed on Friday.</q-item-label>
            </q-item-section>
            <q-item-section side top>
                <q-item-label caption>Yesterday</q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
  </q-page>
</template>
