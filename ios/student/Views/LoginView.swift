import SwiftUI

struct LoginView: View {
    @Binding var isLoggedIn: Bool
    @Binding var token: String
    @Binding var studentName: String
    var apiService: StudentAPIServiceProtocol = HttpStudentAPIService()
    
    @State private var email = ""
    @State private var password = ""
    @State private var isLoading = false
    @State private var errorMessage: String? = nil

    var body: some View {
        ZStack {
            StudentTheme.primaryGradient
                .ignoresSafeArea()
            
            VStack(spacing: 24) {
                VStack(spacing: 12) {
                    Image(systemName: "graduationcap.fill")
                        .font(.system(size: 60))
                        .foregroundColor(.white)
                    
                    Text("Registro Studente")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                        .foregroundColor(.white)
                    
                    Text("Accedi al tuo account")
                        .font(.subheadline)
                        .foregroundColor(.white.opacity(0.8))
                }
                
                VStack(spacing: 16) {
                    HStack {
                        Image(systemName: "envelope.fill")
                            .foregroundColor(.gray)
                        TextField("Email Studente", text: $email)
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
                                    self.studentName = "Studente"
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
                            Text("ACCEDI")
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
                .background(Color(UIColor.secondarySystemGroupedBackground).opacity(0.95))
                .cornerRadius(24)
                .padding(.horizontal, 20)
            }
        }
    }
}
