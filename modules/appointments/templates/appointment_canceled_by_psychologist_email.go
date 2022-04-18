package appointments_templates

// AppointmentCanceledByPsychologistEmailTemplate is an email template used to tell the psychologist when a appointment is canceled by the psychologist
var AppointmentCanceledByPsychologistEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que {{ .PsyFullName }} cancelou a sua próxima consulta.</p>
<p>Entre no nosso site caso queira sugerir um novo horário.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`
