import SwiftUI

struct StudentELearningView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Quiz Formativi & E-Learning")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Quiz: Limiti e Continuità")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Aperto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Matematica • 10 Domande (20 min)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Tentativi: 2 • Punteggio minimo: 60%")
                            .font(.caption)
                            .foregroundColor(.secondary)

                        Button("Avvia Quiz") {}
                            .buttonStyle(.borderedProminent)
                            .padding(.top, 4)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("E-Learning & Quiz")
        }
    }
}
