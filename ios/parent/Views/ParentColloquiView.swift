import SwiftUI

struct ParentColloquiView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Docenti Disponibili per Ricevimento")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Prof.ssa Bianchi (Matematica)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Disponibile")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Mercoledì • 10:00 - 11:00 (Aula 12)")
                            .font(.subheadline)
                            .foregroundColor(.blue)

                        Button("Prenota Slot") {}
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Colloqui Settimanali")
        }
    }
}
