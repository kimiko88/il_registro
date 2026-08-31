import SwiftUI

struct TeacherCoTeachingView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Sessioni di Compresenza")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Scienze + Lab Chimico")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Doppia Firma")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Prof.ssa Bianchi (Teoria) + Prof. Neri (ITP)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Classe 3ª A • 3ª Ora • Analisi composti organici")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Compresenza")
        }
    }
}
