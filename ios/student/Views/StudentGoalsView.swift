import SwiftUI

struct StudentGoalsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Traguardi di Miglioramento")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Studio di Funzione")
                                .fontWeight(.bold)
                            Spacer()
                            Text("75%")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("Matematica • Target: Media 7.0")
                            .font(.subheadline)
                            .foregroundColor(.blue)

                        ProgressView(value: 0.75)
                            .padding(.vertical, 2)

                        Text("Feedback: Ottimo impegno nei compiti a casa")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Obiettivi Didattici")
        }
    }
}
