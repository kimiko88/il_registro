<template>
  <q-page padding class="bg-slate-50">
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <h1 class="text-h4 text-weight-bold text-slate-800 q-my-none">Gestione Classi</h1>
        <p class="text-subtitle1 text-slate-500 q-mb-none">Pianificazione classi, cattedre e adozioni libri</p>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-select
          v-model="selectedYear"
          :options="academicYearOptions"
          label="Anno Accademico"
          outlined
          dense
          class="rounded-lg min-width-150"
          @update:model-value="onYearChange"
        />
        <q-btn
          color="indigo-7"
          icon="published_with_changes"
          label="Migrazione Anno"
          class="rounded-lg shadow-sm"
          no-caps
          @click="openMigrationWizard"
        />
        <q-btn color="primary" icon="add" label="Nuova Classe" class="rounded-lg shadow-sm" @click="openDialog()" />
      </div>
    </div>

    <!-- Classes List -->
    <q-card class="rounded-xl shadow-soft border-slate-100 overflow-hidden bg-white">
      <q-table
        :rows="classesStore.classes"
        :columns="columns"
        :filter="filter"
        :loading="classesStore.loading"
        row-key="id"
        flat
        class="bg-transparent"
        :pagination="{ rowsPerPage: 10 }"
      >
        <template #top-right>
          <q-input dense debounce="300" v-model="filter" placeholder="Cerca classe..." outlined>
            <template #append>
              <q-icon name="search" color="grey-5" />
            </template>
          </q-input>
        </template>
        
        <template #header-cell="props">
          <q-th :props="props" class="text-slate-500 font-bold">
            {{ props.col.label }}
          </q-th>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-right">
            <q-btn flat round dense icon="menu_book" color="indigo" @click="openAssignmentsDialog(props.row)">
              <q-tooltip>Gestione Materie & Docenti</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="groups" color="cyan-9" @click="openStudentsDialog(props.row)">
              <q-tooltip>Gestione Studenti della Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="auto_stories" color="emerald" @click="openTextbooksDialog(props.row)">
              <q-tooltip>Libri di Testo</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="calendar_today" color="orange" @click="openScheduleDialog(props.row)">
              <q-tooltip>Orario Settimanale</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="edit" color="primary" @click="openDialog(props.row)">
              <q-tooltip>Modifica Classe</q-tooltip>
            </q-btn>
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)">
              <q-tooltip>Elimina Classe</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Create/Edit Class -->
    <q-dialog v-model="showDialog" persistent class="premium-dialog">
      <q-card style="width: min(480px, 95vw); max-width: 95vw;" class="rounded-xl overflow-hidden shadow-24">
        <q-card-section class="bg-gradient-primary text-white row items-center q-pa-lg">
          <div class="text-h6 text-weight-bold">{{ isEdit ? 'Modifica Classe' : 'Nuova Classe' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup :aria-label="$t('common.close') || 'Chiudi'" />
        </q-card-section>

        <q-card-section class="q-pa-xl">
          <q-form @submit="saveClass" class="q-gutter-y-lg">
            <div class="row q-col-gutter-lg">
              <div class="col-8">
                <q-input v-model="form.name" label="Nome (es. 1, 5)" outlined :rules="[val => !!val || 'Obbligatorio']" />
              </div>
              <div class="col-4">
                <q-input v-model="form.section" label="Sezione (es. A, B)" outlined :rules="[val => !!val || 'Obbligatorio']" />
              </div>
            </div>
            
            <q-input v-model="form.articolazione" label="Articolazione (es. Informatica, Telecomunicazioni - Opzionale)" outlined />
            <q-input v-model="form.location" label="Sede della Scuola (es. Sede Centrale, Succursale - Opzionale)" outlined />
            
            <q-select
              v-model="form.academic_year"
              :options="academicYearOptions"
              label="Anno Accademico"
              outlined
              :rules="[val => !!val || 'Obbligatorio']"
            />
            
            <q-select
              v-model="form.coordinator_id"
              :options="teacherUserOptions"
              label="Coordinatore di Classe"
              outlined
              emit-value
              map-options
              clearable
            />
            
            <div class="row justify-end q-mt-xl q-gutter-sm">
              <q-btn label="Annulla" flat v-close-popup color="slate-400" />
              <q-btn :label="isEdit ? 'Aggiorna' : 'Crea Classe'" type="submit" color="primary" class="q-px-xl rounded-lg shadow-sm" :loading="saving" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Assignments Dialog -->
    <q-dialog v-model="showAssignmentsDialog">
      <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-gradient-primary text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="menu_book" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Cattedre - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 opacity-80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-7">
              <q-table
                title="Programmazione Didattica"
                :rows="assignments"
                :columns="assignmentsColumns"
                row-key="id"
                flat
                class="bg-transparent border-slate-100 rounded-xl"
              >
                <template #header-cell="props">
                  <q-th :props="props" class="text-slate-500 font-bold">
                    {{ props.col.label }}
                  </q-th>
                </template>

                <template #body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeAssignment(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">Assegna Materia</div>
                <q-form @submit="addAssignment" class="q-gutter-y-md">
                  <q-select
                    v-model="assignForm.subject_id"
                    :options="subjectOptions"
                    label="Materia *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona materia']"
                  >
                    <template #no-option>
                      <q-item>
                        <q-item-section class="text-grey">Nessuna materia trovata</q-item-section>
                      </q-item>
                      <q-item clickable @click="openCreateSubject">
                        <q-item-section class="text-primary text-weight-bold">AGGIUNGI NUOVA MATERIA</q-item-section>
                      </q-item>
                    </template>
                  </q-select>

                  <q-select
                    v-model="assignForm.teacher_id"
                    :options="teacherOptions"
                    label="Docente"
                    outlined
                    dense
                    emit-value
                    map-options
                  />

                  <q-input
                    v-model.number="assignForm.hours_per_week"
                    label="Ore Settimanali"
                    type="number"
                    outlined
                    dense
                    min="1"
                  />

                  <q-btn type="submit" label="Assegna Cattedra" color="primary" class="full-width rounded-lg q-py-sm shadow-sm q-mt-md" no-caps />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Textbooks Dialog -->
    <q-dialog v-model="showTextbooksDialog">
      <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-gradient-premium text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="auto_stories" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Adozioni Libri - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 opacity-80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-7">
              <q-table
                title="Libri Adottati"
                :rows="classTextbooks"
                :columns="textbookColumns"
                row-key="id"
                flat
                class="bg-transparent border-slate-100 rounded-xl"
              >
                <template #header-cell="props">
                  <q-th :props="props" class="text-slate-500 font-bold">
                    {{ props.col.label }}
                  </q-th>
                </template>

                <template #body-cell-actions="props">
                  <q-td :props="props" auto-width>
                    <q-btn flat round dense color="negative" icon="delete" @click="removeTextbook(props.row)" />
                  </q-td>
                </template>
              </q-table>
            </div>
            
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border-slate-200">
                <div class="row items-center justify-between q-mb-md">
                  <div class="text-subtitle1 text-weight-bold text-slate-800">Adotta Libro</div>
                  <q-btn flat dense icon="add" label="Nuovo Libro" color="primary" no-caps @click="openCreateTextbook" />
                </div>
                <q-form @submit="addTextbookToClass" class="q-gutter-y-md">
                  <q-select
                    v-model="textbookForm.textbook_id"
                    :options="allTextbooksOptions"
                    label="Libro *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona libro']"
                  >
                    <template #no-option>
                      <q-item>
                        <q-item-section class="text-grey">Nessun libro in catalogo</q-item-section>
                      </q-item>
                      <q-item clickable @click="openCreateTextbook">
                        <q-item-section class="text-primary text-weight-bold">CREA NUOVO LIBRO</q-item-section>
                      </q-item>
                    </template>
                  </q-select>
                  <q-select
                    v-model="textbookForm.subject_id"
                    :options="subjectOptions"
                    label="Materia *"
                    outlined
                    dense
                    emit-value
                    map-options
                    :rules="[val => !!val || 'Seleziona materia']"
                  />
                  <q-checkbox v-model="textbookForm.is_optional" label="Il libro è opzionale" class="text-slate-700" />
                  <q-btn type="submit" label="Conferma Adozione" color="primary" class="full-width rounded-lg q-py-sm shadow-sm q-mt-md" no-caps />
                </q-form>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Schedule Dialog -->
    <q-dialog v-model="showScheduleDialog">
      <q-card style="width: min(1200px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-primary text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="calendar_today" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Orario Settimanale - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 text-white/80">{{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <ScheduleGrid 
            :assignments="assignments"
            :initial-schedule="currentSchedule"
            :loading="scheduleLoading"
            @save="saveSchedule"
          />
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Students Assignment Dialog -->
    <q-dialog v-model="showStudentsDialog">
      <q-card style="width: min(1100px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <q-card-section class="bg-cyan-8 text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="groups" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Gestione Studenti - Classe {{ currentClass?.name }}{{ currentClass?.section }}</div>
              <div class="text-subtitle2 text-white/90">{{ classStudents.length }} studenti iscritti · Anno Accademico {{ currentClass?.academic_year }}</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md col overflow-y-auto">
          <div v-if="loadingStudents" class="text-center q-pa-xl">
            <q-spinner-dots color="cyan-8" size="40px" />
          </div>

          <div v-else class="row q-col-gutter-md">
            <!-- Left Column: Currently Enrolled Students -->
            <div class="col-12 col-md-7">
              <div class="row items-center justify-between q-mb-sm">
                <div class="text-subtitle1 text-weight-bold text-slate-800">
                  Studenti in Classe ({{ classStudents.length }})
                </div>
                <q-input
                  v-model="studentSearchFilter"
                  placeholder="Filtra studente..."
                  dense outlined
                  class="bg-white rounded-lg min-width-200"
                >
                  <template #append>
                    <q-icon name="search" size="xs" />
                  </template>
                </q-input>
              </div>

              <div v-if="filteredClassStudents.length === 0" class="bg-slate-50 border border-slate-200 rounded-xl p-6 text-center text-slate-400">
                <q-icon name="person_off" size="48px" class="q-mb-xs opacity-40" />
                <div class="text-subtitle2">Nessuno studente in questa classe</div>
                <div class="text-caption">Usa il pannello a destra per associare gli studenti della scuola.</div>
              </div>

              <q-scroll-area style="height: 440px;" class="rounded-xl border border-slate-200 bg-white">
                <q-list separator>
                  <q-item v-for="student in filteredClassStudents" :key="student.id" class="q-py-sm items-center justify-between">
                    <q-item-section avatar>
                      <q-avatar color="cyan-1" text-color="cyan-9" icon="person" font-size="20px" />
                    </q-item-section>

                    <q-item-section class="col overflow-hidden q-pr-sm">
                      <q-item-label class="text-weight-bold text-slate-800 ellipsis">
                        {{ student.last_name || '' }} {{ student.first_name || student.name || '' }}
                      </q-item-label>
                      <q-item-label caption class="text-slate-500 ellipsis">
                        {{ student.email }} {{ student.enrollment_number ? '· Matr: ' + student.enrollment_number : '' }}
                      </q-item-label>
                    </q-item-section>

                    <q-item-section side class="shrink-0">
                      <q-btn flat round dense color="negative" icon="person_remove" @click="removeStudentFromClass(student)">
                        <q-tooltip>Rimuovi dalla classe</q-tooltip>
                      </q-btn>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-scroll-area>
            </div>

            <!-- Right Column: Add / Transfer Student -->
            <div class="col-12 col-md-5">
              <q-card flat class="rounded-xl bg-slate-50 q-pa-md border border-slate-200">
                <div class="text-subtitle1 text-weight-bold text-slate-800 q-mb-md">
                  Associa o Trasferisci Studente
                </div>

                <div class="q-gutter-y-md">
                  <q-select
                    v-model="selectedStudentToAssign"
                    :options="assignableStudentOptions"
                    label="Seleziona Studente *"
                    outlined dense
                    emit-value map-options
                    filterable
                  >
                    <template #no-option>
                      <q-item>
                        <q-item-section class="text-grey">Tutti gli studenti sono già in questa classe</q-item-section>
                      </q-item>
                    </template>
                  </q-select>

                  <q-btn
                    color="cyan-8"
                    label="Associa alla Classe"
                    icon="person_add"
                    class="full-width rounded-lg q-py-sm shadow-sm"
                    no-caps
                    :disabled="!selectedStudentToAssign"
                    @click="assignStudentToClass(selectedStudentToAssign)"
                  />
                </div>

                <q-separator class="q-my-md" />

                <!-- Unassigned Students Quick Pick -->
                <div class="text-subtitle2 text-weight-bold text-slate-700 q-mb-xs">
                  Studenti Senza Classe ({{ unassignedStudents.length }})
                </div>
                <div v-if="unassignedStudents.length === 0" class="text-caption text-slate-400">
                  Tutti gli studenti iscritti risultano assegnati ad una classe.
                </div>
                <q-list v-else separator class="rounded-lg border border-slate-200 bg-white overflow-y-auto overflow-x-hidden" style="max-height: 200px;">
                  <q-item v-for="uSt in unassignedStudents" :key="uSt.id" class="q-py-xs items-center justify-between">
                    <q-item-section class="col overflow-hidden q-pr-xs">
                      <q-item-label class="text-caption text-weight-bold text-slate-800 ellipsis">
                        {{ uSt.last_name }} {{ uSt.first_name || uSt.name }}
                      </q-item-label>
                      <q-item-label caption class="text-slate-400 text-xs ellipsis">
                        {{ uSt.email }}
                      </q-item-label>
                    </q-item-section>
                    <q-item-section side class="shrink-0">
                      <q-btn flat round size="sm" color="cyan-8" icon="person_add" @click="assignStudentToClass(uSt.id)">
                        <q-tooltip>Associa alla classe</q-tooltip>
                      </q-btn>
                    </q-item-section>
                  </q-item>
                </q-list>
              </q-card>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- Quick Create Subject Dialog -->
    <q-dialog v-model="showSubjectDialog">
      <q-card style="min-width: 350px" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="q-pa-lg">
          <div class="text-h6 text-weight-bold text-slate-800">Nuova Materia</div>
        </q-card-section>
        <q-card-section class="q-px-lg q-pb-lg">
          <q-input v-model="newSubjectName" label="Nome Materia" outlined autofocus @keyup.enter="createSubject" />
        </q-card-section>
        <q-card-actions align="right" class="q-pa-lg bg-slate-50">
          <q-btn flat label="Annulla" v-close-popup color="slate-400" />
          <q-btn label="Crea" color="primary" class="rounded-lg q-px-lg" @click="createSubject" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Create Textbook Dialog -->
    <q-dialog v-model="showCreateTextbookDialog">
      <q-card style="display: flex; flex-direction: column; width: min(550px, 95vw); max-height: 90vh;" class="rounded-xl shadow-24 bg-white">
        <q-card-section class="bg-primary text-white row items-center justify-between q-pa-md">
          <div class="text-subtitle1 font-bold">
            <q-icon name="menu_book" class="q-mr-xs" />
            Nuovo Libro Scolastico (Catalogo)
          </div>
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pa-md scroll" style="flex: 1; overflow-y: auto;">

          <q-form @submit="createTextbookInCatalog" class="q-gutter-y-md">
            <q-input v-model="newTextbook.title" label="Titolo del Libro *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.author" label="Autore / Autori *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.publisher" label="Casa Editrice *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            <q-input v-model="newTextbook.isbn" label="Codice ISBN *" outlined dense :rules="[val => !!val || 'Obbligatorio']" />
            
            <q-select
              v-model="newTextbook.subject_id"
              :options="subjectOptions"
              label="Materia associata"
              outlined dense emit-value map-options clearable
            />
            
            <div class="row q-col-gutter-md">
              <div class="col-6">
                <q-input v-model.number="newTextbook.price" label="Prezzo (€)" type="number" step="0.01" outlined dense />
              </div>
              <div class="col-6">
                <q-input v-model="newTextbook.volume" label="Volume (es. 1, Unico)" outlined dense />
              </div>
            </div>

            <div class="row justify-end q-gutter-sm q-mt-md">
              <q-btn flat label="Annulla" v-close-popup no-caps />
              <q-btn type="submit" label="Salva Libro" color="primary" class="rounded-lg q-px-md" no-caps :loading="savingTextbook" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>


    <!-- Academic Year Migration Wizard Dialog -->
    <q-dialog v-model="showMigrationDialog" persistent class="premium-dialog">
      <q-card style="width: min(1000px, 95vw); max-height: 90vh;" class="rounded-xl overflow-hidden shadow-24 bg-white column no-wrap">
        <!-- Header -->
        <q-card-section class="bg-indigo-9 text-white row items-center q-pa-md shrink-0">
          <div class="row items-center">
            <q-avatar color="white-20" text-color="white" icon="published_with_changes" class="q-mr-sm" size="36px" />
            <div>
              <div class="text-h6 text-weight-bold">Migrazione Studenti e Passaggio d'Anno</div>
              <div class="text-subtitle2 text-white/90">Gestisci avanzamento classi, promossi, bocciati e diplomati</div>
            </div>
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <!-- Body / Stepper -->
        <q-card-section class="q-pa-md col overflow-y-auto">
          <!-- Step 1: Configuration -->
          <div v-if="migrationStep === 1" class="q-gutter-y-md">
            <div class="bg-indigo-50 border border-indigo-100 rounded-xl q-pa-md text-indigo-9">
              <div class="text-weight-bold flex items-center gap-2">
                <q-icon name="info" size="sm" /> Procedura di Passaggio d'Anno
              </div>
              <div class="text-caption q-mt-xs">
                Seleziona l'Anno Scolastico sorgente di origine (es. 2024/2025) e quello di destinazione per il nuovo anno.
                Potrai configurare rapidamente per ciascuno studente lo stato (Promosso, Bocciato, Diplomato o Trasferito).
              </div>
            </div>

            <div class="row q-col-gutter-md q-mt-sm">
              <div class="col-12 col-md-6">
                <q-select
                  v-model="migrationSourceYear"
                  :options="academicYearOptions"
                  label="Anno Sorgente (Origine)"
                  outlined dense
                />
              </div>
              <div class="col-12 col-md-6">
                <q-select
                  v-model="migrationTargetYear"
                  :options="academicYearOptions"
                  label="Anno Destinazione (Nuovo Anno)"
                  outlined dense
                />
              </div>
              <div class="col-12">
                <q-select
                  v-model="migrationSourceClassId"
                  :options="migrationClassOptions"
                  label="Filtra per Classe (oppure Tutte le Classi)"
                  outlined dense
                  emit-value map-options
                />
              </div>
            </div>
          </div>

          <!-- Step 2: Student Outcomes & Classes Setup -->
          <div v-else-if="migrationStep === 2" class="q-gutter-y-md">
            <!-- Group Actions Toolbar -->
            <div class="row items-center justify-between bg-slate-50 border border-slate-200 rounded-xl q-pa-sm">
              <div class="text-subtitle2 text-weight-bold text-slate-700">
                Studenti da Elaborare: {{ migrationStudents.length }}
              </div>
              <div class="row q-gutter-xs">
                <q-btn size="xs" color="positive" icon="done_all" label="Segna Tutti Promossi" no-caps unelevated @click="setAllMigrationAction('promoted')" />
                <q-btn size="xs" color="purple" icon="school" label="Diploma 5° / Promuovi 1-4°" no-caps unelevated @click="setSmartMigrationDefaults()" />
                <q-btn size="xs" color="negative" icon="block" label="Segna Tutti Bocciati" no-caps unelevated @click="setAllMigrationAction('repeater')" />
              </div>
            </div>

            <!-- Loading dots if fetching -->
            <div v-if="migrationLoading" class="text-center q-pa-xl">
              <q-spinner-dots color="indigo-7" size="40px" />
              <div class="text-caption text-slate-500 q-mt-sm">Caricamento studenti e predisposizione classi dell'anno destinazione...</div>
            </div>

            <!-- Students List with Actions -->
            <div v-else-if="migrationStudents.length === 0" class="text-center q-pa-xl text-slate-400 border rounded-xl bg-slate-50">
              Nessuno studente trovato per i criteri selezionati.
            </div>

            <q-scroll-area v-else style="height: 420px;" class="rounded-xl border border-slate-200 bg-white">
              <q-list separator dense>
                <q-item v-for="st in migrationStudents" :key="st.student_id" class="q-py-sm items-center justify-between">
                  <!-- Student Info -->
                  <q-item-section style="width: 250px" class="shrink-0">
                    <q-item-label class="text-weight-bold text-slate-800 ellipsis">
                      {{ st.last_name }} {{ st.first_name || st.name }}
                    </q-item-label>
                    <q-item-label caption class="text-slate-400 text-xs ellipsis">
                      Classe Attuale: <q-badge color="cyan-8" :label="st.current_class_name" class="q-ml-xs" />
                    </q-item-label>
                  </q-item-section>

                  <!-- Action Buttons -->
                  <q-item-section class="col q-px-sm">
                    <q-btn-toggle
                      v-model="st.action"
                      dense
                      toggle-color="indigo-7"
                      size="xs"
                      no-caps
                      spread
                      :options="[
                        { label: '🟢 Promosso/a', value: 'promoted' },
                        { label: '🔴 Bocciato/a', value: 'repeater' },
                        { label: '🎓 Diplomato/a', value: 'graduated' },
                        { label: '🚪 Trasferito/a', value: 'left' }
                      ]"
                      @update:model-value="onStudentActionChange(st)"
                    />
                  </q-item-section>

                  <!-- Target Class Selector -->
                  <q-item-section style="width: 220px" class="shrink-0 q-pl-sm">
                    <q-select
                      v-if="st.action === 'promoted' || st.action === 'repeater'"
                      v-model="st.target_class_id"
                      :options="getTargetClassOptions(st)"
                      label="Classe Destinazione"
                      outlined dense emit-value map-options
                      style="font-size: 11px;"
                    />
                    <div v-else class="text-caption text-slate-400 text-italic text-center">
                      Nessuna classe (Disassociato)
                    </div>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-scroll-area>
          </div>

          <!-- Step 3: Confirmation Summary -->
          <div v-else-if="migrationStep === 3" class="q-gutter-y-md">
            <div class="text-subtitle1 text-weight-bold text-slate-800">
              Riepilogo Migrazione (da {{ migrationSourceYear }} a {{ migrationTargetYear }})
            </div>

            <div class="row q-col-gutter-md">
              <div class="col-6 col-md-3">
                <q-card flat class="bg-emerald-50 border border-emerald-200 text-emerald-9 q-pa-md text-center rounded-xl">
                  <div class="text-h4 text-weight-bold">{{ migrationSummary.promoted }}</div>
                  <div class="text-caption text-weight-medium">Promossi</div>
                </q-card>
              </div>
              <div class="col-6 col-md-3">
                <q-card flat class="bg-rose-50 border border-rose-200 text-rose-9 q-pa-md text-center rounded-xl">
                  <div class="text-h4 text-weight-bold">{{ migrationSummary.repeater }}</div>
                  <div class="text-caption text-weight-medium">Bocciati</div>
                </q-card>
              </div>
              <div class="col-6 col-md-3">
                <q-card flat class="bg-purple-50 border border-purple-200 text-purple-9 q-pa-md text-center rounded-xl">
                  <div class="text-h4 text-weight-bold">{{ migrationSummary.graduated }}</div>
                  <div class="text-caption text-weight-medium">Diplomati</div>
                </q-card>
              </div>
              <div class="col-6 col-md-3">
                <q-card flat class="bg-slate-100 border border-slate-200 text-slate-7 q-pa-md text-center rounded-xl">
                  <div class="text-h4 text-weight-bold">{{ migrationSummary.left }}</div>
                  <div class="text-caption text-weight-medium">Trasferiti / Ritirati</div>
                </q-card>
              </div>
            </div>
          </div>
        </q-card-section>

        <!-- Footer Actions -->
        <q-card-actions align="right" class="q-pa-md bg-slate-50 border-t border-slate-200 shrink-0">
          <q-btn v-if="migrationStep > 1" flat label="Indietro" color="slate-600" :disabled="migrationLoading" @click="migrationStep--" />
          <q-space />
          <q-btn flat label="Annulla" v-close-popup color="slate-500" />
          <q-btn
            v-if="migrationStep === 1"
            color="indigo-7"
            label="Avanti: Configura Studenti"
            icon-right="arrow_forward"
            no-caps class="rounded-lg q-px-md"
            @click="goToStep2"
          />
          <q-btn
            v-else-if="migrationStep === 2"
            color="indigo-7"
            label="Avanti: Verifica Riepilogo"
            icon-right="arrow_forward"
            no-caps class="rounded-lg q-px-md"
            :disabled="migrationStudents.length === 0"
            @click="migrationStep = 3"
          />
          <q-btn
            v-else-if="migrationStep === 3"
            color="emerald-7"
            label="Conferma ed Esegui Migrazione"
            icon="check_circle"
            no-caps class="rounded-lg q-px-lg shadow-sm"
            :loading="migrationLoading"
            @click="executeMigration"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

  </q-page>
</template>

<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { useClassesStore } from '@/stores/classes'
import { useAuthStore } from '@/stores/auth'
import adminService from '@/services/adminService'
import api from '@/services/api'
import { useQuasar } from 'quasar'
import { textbookService } from '@/services/textbookService'
import { useI18n } from 'vue-i18n'
import ScheduleGrid from '@/components/Secretary/ScheduleGrid.vue'

const $q = useQuasar()
const { t } = useI18n()
const classesStore = useClassesStore()
const authStore = useAuthStore()

const filter = ref('')
const showDialog = ref(false)
const isEdit = ref(false)
const saving = ref(false)
const getCurrentAcademicYear = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1; // 1-12
  // In Italy, the academic year usually starts in September
  if (month >= 9) { 
    return `${year}/${year + 1}`;
  } else {
    return `${year - 1}/${year}`;
  }
}

const currentYearStr = getCurrentAcademicYear();
const selectedYear = ref(currentYearStr);

const currentStart = parseInt(currentYearStr.split('/')[0]);
const academicYearOptions = [
  `${currentStart - 1}/${currentStart}`,
  currentYearStr,
  `${currentStart + 1}/${currentStart + 2}`
]

// Students Management State
const showStudentsDialog = ref(false)
const classStudents = ref([])
const allSchoolStudents = ref([])
const loadingStudents = ref(false)
const selectedStudentToAssign = ref(null)
const studentSearchFilter = ref('')

const sortAlphabetically = (arr) => {
  return [...arr].sort((a, b) => {
    const lastA = (a.last_name || '').toLowerCase()
    const lastB = (b.last_name || '').toLowerCase()
    if (lastA !== lastB) return lastA.localeCompare(lastB, 'it')
    const firstA = (a.first_name || a.name || '').toLowerCase()
    const firstB = (b.first_name || b.name || '').toLowerCase()
    return firstA.localeCompare(firstB, 'it')
  })
}

const filteredClassStudents = computed(() => {
  let list = classStudents.value
  if (studentSearchFilter.value) {
    const q = studentSearchFilter.value.toLowerCase().trim()
    list = list.filter(s =>
      (s.first_name || s.name || '').toLowerCase().includes(q) ||
      (s.last_name || '').toLowerCase().includes(q) ||
      (s.email || '').toLowerCase().includes(q)
    )
  }
  return sortAlphabetically(list)
})

const unassignedStudents = computed(() => {
  const list = allSchoolStudents.value.filter(s => !s.class_id || s.class_id === '')
  return sortAlphabetically(list)
})

const assignableStudentOptions = computed(() => {
  const sorted = sortAlphabetically(allSchoolStudents.value.filter(s => s.class_id !== currentClass.value?.id))
  return sorted.map(s => {
    const classNameStr = s.class_name ? ` (Attualmente in ${s.class_name})` : ' (Senza classe)'
    return {
      label: `${s.last_name || ''} ${s.first_name || s.name || ''}${classNameStr}`.trim(),
      value: s.id
    }
  })
})

const openStudentsDialog = async (row) => {
  currentClass.value = row
  showStudentsDialog.value = true
  selectedStudentToAssign.value = null
  studentSearchFilter.value = ''
  await fetchStudentsForClass()
}

const fetchStudentsForClass = async () => {
  if (!currentClass.value) return
  loadingStudents.value = true
  try {
    const [classRes, allRes] = await Promise.all([
      api.get('/users', { params: { role: 'student', class_id: currentClass.value.id, page_size: 500 } }),
      api.get('/users', { params: { role: 'student', school_id: authStore.user.school_id, page_size: 500 } })
    ])
    classStudents.value = classRes.data?.users || classRes.data || []
    allSchoolStudents.value = allRes.data?.users || allRes.data || []
  } catch (e) {
    console.error('Failed to load students for class', e)
    $q.notify({ type: 'negative', message: 'Errore durante il caricamento degli studenti' })
  } finally {
    loadingStudents.value = false
  }
}

const assignStudentToClass = async (studentId) => {
  if (!studentId || !currentClass.value) return
  try {
    await adminService.updateUser(studentId, { class_id: currentClass.value.id })
    $q.notify({ type: 'positive', message: 'Studente associato alla classe con successo' })
    selectedStudentToAssign.value = null
    await fetchStudentsForClass()
    refreshClasses()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Errore durante l\'associazione dello studente' })
  }
}

const removeStudentFromClass = async (student) => {
  $q.dialog({
    title: 'Rimuovi dalla Classe',
    message: `Sei sicuro di voler rimuovere ${student.first_name || ''} ${student.last_name || ''} dalla classe ${currentClass.value?.name || ''}${currentClass.value?.section || ''}?`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Rimuovi' }
  }).onOk(async () => {
    try {
      await adminService.updateUser(student.id, { class_id: '' })
      $q.notify({ type: 'positive', message: 'Studente rimosso dalla classe' })
      await fetchStudentsForClass()
      refreshClasses()
    } catch (e) {
      $q.notify({ type: 'negative', message: 'Errore durante la rimozione dello studente' })
    }
  })
}

// Assignments State
const showAssignmentsDialog = ref(false)
const showSubjectDialog = ref(false)
const currentClass = ref(null)
const assignments = ref([])
const subjects = ref([])
const teachers = ref([])
const newSubjectName = ref('')

const showScheduleDialog = ref(false)
const scheduleLoading = ref(false)
const currentSchedule = ref([])

const assignForm = reactive({
    subject_id: null,
    teacher_id: null,
    hours_per_week: 1
})

const showTextbooksDialog = ref(false)
const showCreateTextbookDialog = ref(false)
const savingTextbook = ref(false)
const classTextbooks = ref([])
const allTextbooks = ref([])
const textbookForm = reactive({
    textbook_id: null,
    subject_id: null,
    is_optional: false
})
const newTextbook = reactive({
    title: '',
    author: '',
    publisher: '',
    isbn: '',
    subject_id: null,
    price: null,
    volume: ''
})

const form = reactive({
  id: null,
  name: '',
  section: '',
  articolazione: '',
  location: '',
  academic_year: currentYearStr,
  coordinator_id: ''
})

const columns = [
  { name: 'name', label: 'Classe', align: 'left', field: row => `${row.name || ''}${row.section || ''}${row.articolazione ? ' - ' + row.articolazione : ''}`, sortable: true },
  { name: 'location', label: 'Sede', align: 'center', field: row => row.location || 'Sede Centrale', sortable: true },
  { name: 'academic_year', label: 'Anno Accademico', align: 'center', field: 'academic_year', sortable: true },
  { name: 'actions', label: 'Azioni', align: 'right' }
]

const assignmentsColumns = [
    { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left', sortable: true },
    { name: 'teacher', label: 'Docente', field: row => row.teacher_name || 'N/A', align: 'left' },
    { name: 'hours', label: 'Ore/Sett', field: 'hours_per_week', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const textbookColumns = [
    { name: 'subject', label: 'Materia', field: 'subject_name', align: 'left' },
    { name: 'title', label: 'Titolo', field: 'title', align: 'left' },
    { name: 'author', label: 'Autore', field: 'author', align: 'left' },
    { name: 'optional', label: 'Opz.', field: row => row.is_optional ? 'Sì' : 'No', align: 'center' },
    { name: 'actions', label: 'Azioni', align: 'right' }
]

const subjectOptions = computed(() => subjects.value.map(s => ({ label: s.name, value: s.id })))
const teacherOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.id })))
const teacherUserOptions = computed(() => teachers.value.map(t => ({ label: `${t.last_name} ${t.first_name}`, value: t.user_id })))
const allTextbooksOptions = computed(() => allTextbooks.value.map(b => ({ label: b.title, value: b.id })))

onMounted(() => {
  if (authStore.user?.school_id) {
    refreshClasses()
    fetchSchoolData()
  }
})

const refreshClasses = () => {
  classesStore.fetchClasses({ 
    school_id: authStore.user.school_id,
    academic_year: selectedYear.value
  })
}

// Watch subject selection to filter teachers
watch(() => assignForm.subject_id, async (newVal) => {
    assignForm.teacher_id = null
    if (newVal) {
        try {
            const res = await adminService.getTeachersList(authStore.user.school_id, newVal)
            teachers.value = res.data || []
        } catch(e) {
            console.error("Error filtering teachers", e)
        }
    } else {
        // Reset to all teachers
        fetchSchoolData()
    }
})

const onYearChange = () => {
  refreshClasses()
}

const fetchSchoolData = async () => {
    try {
        const [sRes, tRes, bRes] = await Promise.all([
            adminService.getSubjects(authStore.user.school_id),
            adminService.getTeachersList(authStore.user.school_id),
            textbookService.getAll()
        ])
        subjects.value = sRes.data || []
        teachers.value = tRes.data || []
        allTextbooks.value = bRes.data || []
    } catch(e) {
        console.error("Error loading school data", e)
    }
}

const openDialog = (row = null) => {
  if (row) {
    isEdit.value = true
    Object.assign(form, row)
  } else {
    isEdit.value = false
    Object.assign(form, {
      id: null,
      name: '',
      section: '',
      articolazione: '',
      location: '',
      academic_year: selectedYear.value,
      coordinator_id: ''
    })
  }
  showDialog.value = true
}

const saveClass = async () => {
  saving.value = true
  try {
    const payload = { ...form, school_id: authStore.user.school_id };
    if (isEdit.value) {
      await classesStore.updateClass(form.id, payload)
      $q.notify({ type: 'positive', message: 'Classe aggiornata' })
    } else {
      await classesStore.createClass(payload)
      $q.notify({ type: 'positive', message: 'Classe creata' })
    }
    showDialog.value = false
    refreshClasses()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Errore nel salvataggio' })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Conferma Eliminazione',
    message: `Vuoi eliminare la classe ${row.name}${row.section}? Tutte le associazioni verranno rimosse.`,
    cancel: true,
    persistent: true,
    ok: { color: 'negative', label: 'Elimina' }
  }).onOk(async () => {
    try {
      await classesStore.deleteClass(row.id)
      $q.notify({ type: 'positive', message: 'Classe eliminata' })
      refreshClasses()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Errore nell\'eliminazione' })
    }
  })
}

// Assignments Logic
const openAssignmentsDialog = async (row) => {
    currentClass.value = row
    showAssignmentsDialog.value = true
    await fetchAssignments(row.id)
}

const fetchAssignments = async (classId) => {
    try {
        const res = await adminService.getClassSubjects(classId)
        assignments.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento materie' })
    }
}

const addAssignment = async () => {
    if (!currentClass.value) return
    try {
        await adminService.assignSubjectToClass(currentClass.value.id, assignForm)
        $q.notify({ type: 'positive', message: 'Materia assegnata' })
        fetchAssignments(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore assegnazione' })
    }
}

const removeAssignment = async (row) => {
    try {
        await adminService.removeClassSubject(row.class_id, row.id)
        $q.notify({ type: 'positive', message: 'Assegnazione rimossa' })
        fetchAssignments(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore rimozione' })
    }
}

const openCreateSubject = () => {
    showSubjectDialog.value = true
}

const createSubject = async () => {
    if (!newSubjectName.value) return
    try {
        const payload = {
            name: newSubjectName.value,
            school_id: authStore.user.school_id
        }
        await adminService.createSubject(payload)
        $q.notify({ type: 'positive', message: 'Materia creata' })
        showSubjectDialog.value = false
        newSubjectName.value = ''
        fetchSchoolData()
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore creazione materia' })
    }
}

// Textbooks Logic
const openTextbooksDialog = async (row) => {
    currentClass.value = row
    showTextbooksDialog.value = true
    fetchClassTextbooks(row.id)
}

const fetchClassTextbooks = async (classId) => {
    try {
        const res = await textbookService.listByClass(classId)
        classTextbooks.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento libri' })
    }
}

const addTextbookToClass = async () => {
    try {
        await textbookService.assignToClass(currentClass.value.id, textbookForm)
        $q.notify({ type: 'positive', message: 'Libro adottato' })
        fetchClassTextbooks(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore adozione libro' })
    }
}

const removeTextbook = async (row) => {
    try {
        await textbookService.removeFromClass(row.id)
        $q.notify({ type: 'positive', message: 'Adozione rimossa' })
        fetchClassTextbooks(currentClass.value.id)
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore rimozione' })
    }
}

const openCreateTextbook = () => {
    Object.assign(newTextbook, {
        title: '',
        author: '',
        publisher: '',
        isbn: '',
        subject_id: textbookForm.subject_id || null,
        price: null,
        volume: ''
    })
    showCreateTextbookDialog.value = true
}

const createTextbookInCatalog = async () => {
    savingTextbook.value = true
    try {
        const payload = {
            ...newTextbook,
            school_id: authStore.user.school_id
        }
        const res = await textbookService.create(payload)
        $q.notify({ type: 'positive', message: 'Nuovo libro creato nel catalogo' })
        showCreateTextbookDialog.value = false
        await fetchSchoolData()
        const createdId = res.data?.id || res.data?.data?.id
        if (createdId) {
            textbookForm.textbook_id = createdId
        }
    } catch (e) {
        console.error('Error creating textbook', e)
        $q.notify({ type: 'negative', message: 'Errore durante la creazione del libro' })
    } finally {
        savingTextbook.value = false
    }
}


// Schedule Logic
const openScheduleDialog = async (row) => {
    currentClass.value = row
    showScheduleDialog.value = true
    scheduleLoading.value = true
    try {
        // We need both assignments (for options) and the current schedule
        await fetchAssignments(row.id)
        const res = await adminService.getClassSchedule(row.id)
        currentSchedule.value = res.data || []
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore caricamento orario' })
    } finally {
        scheduleLoading.value = false
    }
}

const saveSchedule = async (entries) => {
    if (!currentClass.value) return
    scheduleLoading.value = true
    try {
        await adminService.saveClassSchedule(currentClass.value.id, { entries })
        $q.notify({ type: 'positive', message: 'Orario salvato con successo' })
        showScheduleDialog.value = false
    } catch(e) {
        $q.notify({ type: 'negative', message: 'Errore durante il salvataggio dell\'orario' })
    } finally {
        scheduleLoading.value = false
    }
}

// Academic Year Migration State & Methods
const showMigrationDialog = ref(false)
const migrationStep = ref(1)
const migrationLoading = ref(false)
const migrationSourceYear = ref(currentYearStr)
const migrationTargetYear = ref(`${currentStart + 1}/${currentStart + 2}`)
const migrationSourceClassId = ref(null)

const targetYearClasses = ref([])
const migrationStudents = ref([])

const migrationClassOptions = computed(() => {
  return [
    { label: 'Tutte le classi dell\'anno sorgente', value: null },
    ...classesStore.classes.map(c => ({
      label: `${c.name}${c.section} (${c.academic_year})`,
      value: c.id
    }))
  ]
})

const openMigrationWizard = () => {
  migrationStep.value = 1
  migrationSourceYear.value = selectedYear.value || currentYearStr
  migrationSourceClassId.value = null
  showMigrationDialog.value = true
}

const autoCreateMissingTargetClasses = async () => {
  const existingNames = new Set(targetYearClasses.value.map(c => `${c.name}${c.section}`))
  const sourceClasses = migrationSourceClassId.value
    ? classesStore.classes.filter(c => c.id === migrationSourceClassId.value)
    : classesStore.classes

  for (const sc of sourceClasses) {
    const section = sc.section || ''
    const rawName = sc.name || ''
    const match = rawName.match(/^(\d+)(.*)$/)
    const gradeNum = match ? parseInt(match[1]) : null
    const restName = match ? match[2] : rawName

    // Repeater class target
    const repeaterTargetName = `${sc.name}${section}`
    if (!existingNames.has(repeaterTargetName)) {
      try {
        const created = await adminService.createClass({
          name: sc.name,
          section: sc.section,
          articolazione: sc.articolazione || '',
          academic_year: migrationTargetYear.value
        })
        targetYearClasses.value.push(created.data)
        existingNames.add(repeaterTargetName)
      } catch (e) {
        console.error('Failed auto-creating repeater target class', e)
      }
    }

    // Promoted class target
    if (gradeNum && gradeNum < 5) {
      const nextGradeName = `${gradeNum + 1}${restName}`
      const promotedTargetName = `${nextGradeName}${section}`
      if (!existingNames.has(promotedTargetName)) {
        try {
          const created = await adminService.createClass({
            name: nextGradeName,
            section: sc.section,
            articolazione: sc.articolazione || '',
            academic_year: migrationTargetYear.value
          })
          targetYearClasses.value.push(created.data)
          existingNames.add(promotedTargetName)
        } catch (e) {
          console.error('Failed auto-creating promoted target class', e)
        }
      }
    }
  }
}

const goToStep2 = async () => {
  migrationLoading.value = true
  migrationStep.value = 2
  try {
    const targetClassesRes = await adminService.getSchoolClasses(authStore.user.school_id, migrationTargetYear.value)
    targetYearClasses.value = targetClassesRes.data || []

    await autoCreateMissingTargetClasses()

    let sourceClassesToProcess = []
    if (migrationSourceClassId.value) {
      sourceClassesToProcess = classesStore.classes.filter(c => c.id === migrationSourceClassId.value)
    } else {
      sourceClassesToProcess = classesStore.classes.filter(c => c.academic_year === migrationSourceYear.value)
    }

    const studentPromises = sourceClassesToProcess.map(c =>
      api.get('/users', { params: { role: 'student', class_id: c.id, page_size: 500 } }).then(res => ({
        class: c,
        students: res.data?.users || res.data || []
      }))
    )

    const classResults = await Promise.all(studentPromises)
    const items = []

    for (const cr of classResults) {
      const cName = `${cr.class.name}${cr.class.section}`
      const rawName = cr.class.name || ''
      const match = rawName.match(/^(\d+)(.*)$/)
      const gradeNum = match ? parseInt(match[1]) : null
      const restName = match ? match[2] : rawName

      for (const st of cr.students) {
        let action = 'promoted'
        if (gradeNum >= 5) {
          action = 'graduated'
        }

        let targetClassId = null
        if (action === 'promoted' && gradeNum && gradeNum < 5) {
          const nextName = `${gradeNum + 1}${restName}`
          const foundTarget = targetYearClasses.value.find(tc => tc.name === nextName && tc.section === cr.class.section)
          if (foundTarget) targetClassId = foundTarget.id
        }

        items.push({
          student_id: st.id,
          first_name: st.first_name || st.name || '',
          last_name: st.last_name || '',
          email: st.email || '',
          current_class_id: cr.class.id,
          current_class_name: cName,
          current_grade: gradeNum,
          current_section: cr.class.section,
          action: action,
          target_class_id: targetClassId
        })
      }
    }

    migrationStudents.value = sortAlphabetically(items)
  } catch (e) {
    console.error('Error setting up migration step 2', e)
    $q.notify({ type: 'negative', message: 'Errore nel caricamento dati per la migrazione' })
  } finally {
    migrationLoading.value = false
  }
}

const onStudentActionChange = (st) => {
  if (st.action === 'promoted') {
    if (st.current_grade && st.current_grade < 5) {
      const nextName = `${st.current_grade + 1}`
      const foundTarget = targetYearClasses.value.find(tc => tc.name.startsWith(nextName) && tc.section === st.current_section)
      if (foundTarget) st.target_class_id = foundTarget.id
    } else {
      st.target_class_id = null
    }
  } else if (st.action === 'repeater') {
    const foundTarget = targetYearClasses.value.find(tc => tc.name.startsWith(`${st.current_grade}`) && tc.section === st.current_section)
    if (foundTarget) st.target_class_id = foundTarget.id
  } else {
    st.target_class_id = null
  }
}

const setAllMigrationAction = (action) => {
  migrationStudents.value.forEach(st => {
    st.action = action
    onStudentActionChange(st)
  })
}

const setSmartMigrationDefaults = () => {
  migrationStudents.value.forEach(st => {
    if (st.current_grade >= 5) {
      st.action = 'graduated'
    } else {
      st.action = 'promoted'
    }
    onStudentActionChange(st)
  })
}

const getTargetClassOptions = (_st) => {
  return targetYearClasses.value.map(c => ({
    label: `${c.name}${c.section}${c.articolazione ? ' ('+c.articolazione+')' : ''}`,
    value: c.id
  }))
}

const migrationSummary = computed(() => {
  const summary = { promoted: 0, repeater: 0, graduated: 0, left: 0 }
  migrationStudents.value.forEach(st => {
    if (summary[st.action] !== undefined) {
      summary[st.action]++
    }
  })
  return summary
})

const executeMigration = async () => {
  migrationLoading.value = true
  try {
    const payload = {
      source_academic_year: migrationSourceYear.value,
      target_academic_year: migrationTargetYear.value,
      migrations: migrationStudents.value.map(st => ({
        student_id: st.student_id,
        action: st.action,
        target_class_id: (st.action === 'promoted' || st.action === 'repeater') ? st.target_class_id : null
      }))
    }

    await api.post('/classes/migrate-students', payload)

    $q.notify({
      type: 'positive',
      message: 'Migrazione anno scolastico completata con successo!',
      caption: `${migrationSummary.value.promoted} Promossi, ${migrationSummary.value.repeater} Bocciati, ${migrationSummary.value.graduated} Diplomati`
    })

    showMigrationDialog.value = false
    await classesStore.fetchClasses(selectedYear.value)
  } catch (e) {
    console.error('Migration failed', e)
    $q.notify({
      type: 'negative',
      message: 'Errore durante l\'esecuzione della migrazione',
      caption: e.response?.data?.error || e.message
    })
  } finally {
    migrationLoading.value = false
  }
}

defineExpose({
    showDialog,
    isEdit,
    form,
    showAssignmentsDialog,
    showStudentsDialog,
    showSubjectDialog,
    currentClass,
    assignments,
    assignForm,
    newSubjectName,
    openDialog,
    saveClass,
    confirmDelete,
    openAssignmentsDialog,
    openStudentsDialog,
    assignStudentToClass,
    removeStudentFromClass,
    addAssignment,
    removeAssignment,
    openCreateSubject,
    createSubject
})
</script>

<style scoped>
.rounded-xl { border-radius: 1rem; }
.rounded-lg { border-radius: 0.75rem; }
.rounded-md { border-radius: 0.5rem; }
.shadow-soft { box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05); }
.border-slate-100 { border: 1px solid #f1f5f9; }
.bg-slate-50 { background-color: #f8fafc; }
</style>
