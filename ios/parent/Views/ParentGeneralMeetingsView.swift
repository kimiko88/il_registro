import SwiftUI

struct ParentGeneralMeetingsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Fila Virtuale Colloqui Pomeridiani")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Prof.ssa Bianchi (Matematica)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Fila (#3)")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Orario Stimato: 16:15 • Aula 12")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Attesa residua: ~8 minuti")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Ricevimento Generale")
        }
    }
}
