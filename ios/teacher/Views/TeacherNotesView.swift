import SwiftUI

public struct TeacherNotesView: View {
    @State private var student = ""
    @State private var noteType = "Disciplinare"
    @State private var noteText = ""
    @State private var submitted = false
    public var onSubmitNote: ((String, String, String) -> Void)?

    public init(onSubmitNote: ((String, String, String) -> Void)? = nil) {
        self.onSubmitNote = onSubmitNote
    }

    public var body: some View {
        NavigationView {
            Form {
                if submitted {
                    Section {
                        Text("Nota registrata nel database e notifica inviata.")
                            .foregroundColor(.green)
                    }
                } else {
                    Section(header: Text(NSLocalizedString("teacher_dashboard_title", comment: ""))) {
                        TextField("Studente", text: $student)
                        TextField("Tipologia", text: $noteType)
                        TextEditor(text: $noteText)
                            .frame(height: 100)

                        Button(NSLocalizedString("add_grade", comment: "")) {
                            onSubmitNote?(student, noteType, noteText)
                            submitted = true
                        }
                        .disabled(student.isEmpty || noteText.isEmpty)
                    }
                }
            }
            .navigationTitle(NSLocalizedString("teacher_dashboard_title", comment: ""))
        }
    }
}
