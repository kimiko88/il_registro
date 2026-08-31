import SwiftUI

struct AccessibilityStatementView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Conformità AgID / WCAG 2.2 AA")) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Stato: Totalmente Conforme")
                            .fontWeight(.bold)
                            .foregroundColor(.green)
                        Text("Supporto completo a VoiceOver, testo dinamico, alto contrasto e filtri visivi.")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text("Contatti RTD")) {
                    Text("Email RTD: rtd@scuola.edu.it")
                        .font(.caption)
                        .foregroundColor(.blue)
                }
            }
            .navigationTitle("Accessibilità AgID")
        }
    }
}
