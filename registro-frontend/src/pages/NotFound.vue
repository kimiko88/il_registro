<template>
  <q-layout view="LhH Lpr lFf">
    <q-page-container>
      <q-page class="fullscreen flex flex-center q-pa-md bg-slate-50 text-slate-800" role="main">
        <div class="text-center" style="max-width: 500px">
          <!-- Illustration / Icon -->
          <div class="q-mb-md">
            <q-icon name="find_in_page" size="100px" color="primary" class="opacity-80" />
          </div>

          <div class="text-h1 text-weight-bolder text-primary text-outfit q-mb-sm" style="font-size: 5rem; line-height: 1">
            404
          </div>

          <div class="text-h5 text-weight-bold q-mb-sm">
            Pagina non trovata
          </div>

          <p class="text-subtitle1 text-slate-600 q-mb-xl">
            La pagina che stai cercando non esiste o potrebbe essere stata spostata.
          </p>

          <div class="row q-col-gutter-sm justify-center">
            <div class="col-12 col-sm-auto">
              <q-btn
                color="primary"
                unelevated
                icon="dashboard"
                label="Torna alla Dashboard"
                to="/dashboard"
                no-caps
                class="full-width rounded-lg q-px-lg"
              />
            </div>
            <div class="col-12 col-sm-auto">
              <q-btn
                outline
                color="primary"
                icon="arrow_back"
                label="Pagina Precedente"
                @click="$router.back()"
                no-caps
                class="full-width rounded-lg q-px-lg"
              />
            </div>
          </div>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
</script>

<style scoped>
.opacity-80 {
  opacity: 0.8;
}
</style>
