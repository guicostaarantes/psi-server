package users_templates

// ResetPasswordEmailTemplate is an email template used when a user wants to reset their password
var ResetPasswordEmailTemplate = `<h2>Olá 😊</h2>
<p>Você recebeu esse email pois pediu uma mudança de senha no PSI.</p>
<p>Para confirmarmos que este email é realmente seu, clique no botão abaixo e crie uma nova senha pra você acessar nosso site.</p>
<a href="{{ .SiteURL }}/nova-senha?token={{ .Token }}">Mudar minha senha</a>`
