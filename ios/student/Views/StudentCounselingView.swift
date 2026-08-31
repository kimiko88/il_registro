import SwiftUI

struct StudentCounselingView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Sportello di Ascolto CIC")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Colloquio Riservato")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Confermato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Venerdì 29 Maggio • 11:15 - 12:00")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Dott.ssa Elena Moretti • Stanza Ascolto (Ala Ovest)")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Sportello Psicologico")
        }
    }
}
