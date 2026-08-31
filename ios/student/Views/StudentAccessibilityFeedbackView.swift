import SwiftUI

struct StudentAccessibilityFeedbackView: View {
    @State private var barrierType = "Contrasto / Visibilità"
    @State private var description = ""
    @State private var submitted = false

    var body: some View {
        NavigationView {
            Form {
                if submitted {
                    Section {
                        VStack(alignment: .leading, spacing: 6) {
                            Text("Segnalazione Inviata con Successo!")
                                .fontWeight(.bold)
                                .foregroundColor(.green)
                            Text("Protocollo: A11Y-2026-0528-0912")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                    }
                } else {
                    Section(header: Text("Meccanismo di Feedback AgID")) {
                        Text("Segnala all'RTD eventuali barriere di accessibilità o problemi WCAG.")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }

                    Section(header: Text("Dati Segnalazione")) {
                        TextField("Tipologia", text: $barrierType)
                        TextEditor(text: $description)
                            .frame(height: 100)

                        Button("Invia Segnalazione all'RTD") {
                            submitted = true
                        }
                        .disabled(description.isEmpty)
                    }
                }
            }
            .navigationTitle("Accessibilità AgID")
        }
    }
}
