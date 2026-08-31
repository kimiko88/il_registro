import SwiftUI

struct ParentNotesView: View {
    @State private var isAcknowledged = false

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Note Disciplinari e Richiami")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Nota Disciplinare Individuale")
                                .fontWeight(.bold)
                            Spacer()
                            Text(isAcknowledged ? "Firmata" : "Richiesta Firma")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(isAcknowledged ? Color.green.opacity(0.2) : Color.red.opacity(0.2))
                                .foregroundColor(isAcknowledged ? .green : .red)
                                .cornerRadius(4)
                        }
                        Text("Prof. Rossi • 22 Maggio 2026")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Disturbo ripetuto della lezione.")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        if !isAcknowledged {
                            Button("Firma Presa Visione") {
                                isAcknowledged = true
                            }
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 2)
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Note Disciplinari")
        }
    }
}
