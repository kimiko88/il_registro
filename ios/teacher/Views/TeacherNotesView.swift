import SwiftUI

struct TeacherNotesView: View {
    @State private var student = ""
    @State private var noteType = "Disciplinare"
    @State private var noteText = ""
    @State private var submitted = false

    var body: some View {
        NavigationView {
            Form {
                if submitted {
                    Section {
                        Text("Nota registrata e notifica inviata alla famiglia.")
                            .foregroundColor(.green)
                    }
                } else {
                    Section(header: Text("Inserimento Nota Disciplinare")) {
                        TextField("Studente", text: $student)
                        TextField("Tipologia", text: $noteType)
                        TextEditor(text: $noteText)
                            .frame(height: 100)

                        Button("Registra Nota") {
                            submitted = true
                        }
                        .disabled(student.isEmpty || noteText.isEmpty)
                    }
                }
            }
            .navigationTitle("Note Disciplinari")
        }
    }
}
