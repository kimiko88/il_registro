import SwiftUI

struct TeacherVerbaliView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Verbali del Consiglio di Classe")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Verbale N. 5 - PDP & Libri")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Approvato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("18 Maggio 2026 • Presidente: Prof.ssa Verdi")
                            .font(.subheadline)
                            .foregroundColor(.blue)

                        Button(action: {}) {
                            Label("Visualizza Verbale PDF", systemImage: "doc.text.fill")
                        }
                        .buttonStyle(.bordered)
                        .padding(.top, 2)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Verbali Consiglio")
        }
    }
}
