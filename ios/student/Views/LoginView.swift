import SwiftUI

public struct LoginView: View {
    @Binding public var isLoggedIn: Bool
    @Binding public var token: String
    @Binding public var studentName: String
    public var apiService: StudentAPIServiceProtocol
    
    @State private var email = "studentea_1@scuola.it"
    @State private var password = "password"
    @State private var isLoading = false
    @State private var errorMessage: String? = nil

    public init(
        isLoggedIn: Binding<Bool>,
        token: Binding<String>,
        studentName: Binding<String>,
        apiService: StudentAPIServiceProtocol = HttpStudentAPIService()
    ) {
        self._isLoggedIn = isLoggedIn
        self._token = token
        self._studentName = studentName
        self.apiService = apiService
    }

    public var body: some View {
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
                    
                    Text(NSLocalizedString("login_subtitle", comment: ""))
                        .font(.subheadline)
                        .foregroundColor(.white.opacity(0.8))
                }
                
                VStack(spacing: 16) {
                    HStack {
                        Image(systemName: "envelope.fill")
                            .foregroundColor(.gray)
                        TextField(NSLocalizedString("email_label", comment: ""), text: $email)
                            .autocapitalization(.none)
                            .keyboardType(.emailAddress)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    HStack {
                        Image(systemName: "lock.fill")
                            .foregroundColor(.gray)
                        SecureField(NSLocalizedString("password_label", comment: ""), text: $password)
                            .onSubmit(performLogin)
                    }
                    .padding()
                    .background(Color(UIColor.systemBackground))
                    .cornerRadius(12)
                    
                    if let errorMessage = errorMessage {
                        Text(errorMessage)
                            .font(.caption)
                            .foregroundColor(.red)
                    }
                    
                    Button(action: performLogin) {
                        if isLoading {
                            ProgressView()
                                .tint(.white)
                        } else {
                            Text(NSLocalizedString("btn_login", comment: ""))
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

    private func performLogin() {
        if email.isEmpty || password.isEmpty {
            errorMessage = NSLocalizedString("err_enter_credentials", comment: "")
            return
        }
        isLoading = true
        errorMessage = nil
        Task {
            do {
                let receivedToken = try await apiService.login(email: email.trimmingCharacters(in: .whitespacesAndNewlines), password: password)
                await MainActor.run {
                    self.token = receivedToken
                    self.studentName = apiService.lastStudentName ?? "Studente"
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
    }
}
