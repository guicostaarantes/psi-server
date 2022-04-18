package users_templates

// ResetPasswordEmailTemplate is an email template used when a user wants to reset their password
var ResetPasswordEmailTemplate = `<h2>Olá 😊</h2>
<p>Você recebeu esse email pois solicitou uma redefinição de senha no PSI.</p>
<p>Clique no botão abaixo para redefinir sua senha e poder acessar nosso site.</p>
<a href="{{ .SiteURL }}/nova-senha?token={{ .Token }}">Redefinir minha senha</a>`
