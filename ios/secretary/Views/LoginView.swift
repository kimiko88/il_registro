import SwiftUI

struct LoginView: View {
    @Binding var isLoggedIn: Bool
    @Binding var token: String
    var apiService: SecretaryAPIServiceProtocol = HttpSecretaryAPIService()
    
    @State private var email = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String? = nil

    var body: some View {
        ZStack {
            Color(red: 0.35, green: 0.11, blue: 0.53)
                .ignoresSafeArea()
            
            VStack(spacing: 24) {
                VStack(spacing: 12) {
                    Image(systemName: "building.columns.fill")
                        .font(.system(size: 60))
                        .foregroundColor(.purple)
                    
                    Text("Registro Segreteria")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                        .foregroundColor(.white)
                    
                    Text("Portale Amministrativo & Didattico")
                        .font(.subheadline)
                        .foregroundColor(.white.opacity(0.8))
                }
                
                VStack(spacing: 16) {
                    HStack {
                        Image(systemName: "envelope.fill")
                            .foregroundColor(.gray)
                        TextField("Email Segreteria", text: $email)
                            .autocapitalization(.none)
                            .keyboardType(.emailAddress)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    HStack {
                        Image(systemName: "lock.fill")
                            .foregroundColor(.gray)
                        SecureField("Password", text: $password)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    if let errorMessage = errorMessage {
                        Text(errorMessage)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                    
                    Button(action: {
                        if email.isEmpty || password.isEmpty {
                            errorMessage = "Inserisci email e password"
                            return
                        }
                        isLoading = true
                        errorMessage = nil
                        Task {
                            do {
                                let receivedToken = try await apiService.login(email: email.trimmingCharacters(in: .whitespacesAndNewlines), password: password)
                                await MainActor.run {
                                    self.token = receivedToken
                                    self.isLoggedIn = true
                                    self.isLoading = false
                                }
                            } catch {
                                await MainActor.run {
                                    self.errorMessage = error.localizedDescription
                                    self.isLoading = false
                                }
                            }
                        }
                    }) {
                        if isLoading {
                            ProgressView()
                                .tint(.white)
                        } else {
                            Text("ACCEDI AL PANNELLO")
                                .font(.headline)
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background(Color.purple)
                                .cornerRadius(12)
                        }
                    }
                    .disabled(isLoading)
                }
                .padding(24)
                .background(Color(UIColor.secondarySystemGroupedBackground))
                .cornerRadius(24)
                .padding(.horizontal, 20)
            }
        }
    }
}
