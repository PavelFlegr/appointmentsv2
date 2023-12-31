// Code generated by templ@v0.2.364 DO NOT EDIT.

package appointment

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import appointmentService "appointmentsv2/services/appointment"
import "appointmentsv2/templates"
import "fmt"

func Reserve(groups []appointmentService.SlotGroup) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_2 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"ui container\">")
			if err != nil {
				return err
			}
			for _, group := range groups {
				_, err = templBuffer.WriteString("<h2 class=\"ui header center aligned\">")
				if err != nil {
					return err
				}
				var var_3 string = group.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_3))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</h2> <div class=\"ui grid\">")
				if err != nil {
					return err
				}
				for _, slot := range *group.Slots {
					_, err = templBuffer.WriteString("<div class=\"ui four wide column\"><div class=\"ui attached segment\"><div class=\"ui attached label\">")
					if err != nil {
						return err
					}
					var_4 := `Spots Left: `
					_, err = templBuffer.WriteString(var_4)
					if err != nil {
						return err
					}
					var var_5 string = fmt.Sprint(slot.Free)
					_, err = templBuffer.WriteString(templ.EscapeString(var_5))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</div></div><div class=\"ui attached segment center aligned\"><div class=\"ui header\">")
					if err != nil {
						return err
					}
					var var_6 string = slot.Start.Local().Format("15:04")
					_, err = templBuffer.WriteString(templ.EscapeString(var_6))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var_7 := `- `
					_, err = templBuffer.WriteString(var_7)
					if err != nil {
						return err
					}
					var var_8 string = slot.End.Local().Format("15:04")
					_, err = templBuffer.WriteString(templ.EscapeString(var_8))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</div><div class=\"ui\"><a class=\"ui primary button\" href=\"")
					if err != nil {
						return err
					}
					var var_9 templ.SafeURL = templ.URL(fmt.Sprintf("/reserve/slot/%v", slot.Token))
					_, err = templBuffer.WriteString(templ.EscapeString(string(var_9)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var_10 := `Reserve`
					_, err = templBuffer.WriteString(var_10)
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</a></div></div></div>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = templates.Layout().Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func ReserveSlot(slot appointmentService.Slot, appointmentToken string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_11 := templ.GetChildren(ctx)
		if var_11 == nil {
			var_11 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		var_12 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<div class=\"ui container\"><div><a href=\"")
			if err != nil {
				return err
			}
			var var_13 templ.SafeURL = templ.URL(fmt.Sprintf("/reserve/%v", appointmentToken))
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_13)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			var_14 := `Back to list`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</a></div><h1>")
			if err != nil {
				return err
			}
			var_15 := `Finish Reservation`
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1><form class=\"ui form\" hx-post=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/reserve/slot/%v", slot.Token)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-vals=\"js:timezone:Intl.DateTimeFormat().resolvedOptions().timeZone\"><div class=\"two fields\"><div class=\"field\"><label>")
			if err != nil {
				return err
			}
			var_16 := `Date`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(slot.Start.Local().Format("2. 1. 2006")))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" disabled></div><div class=\"field\"><label>")
			if err != nil {
				return err
			}
			var_17 := `Time`
			_, err = templBuffer.WriteString(var_17)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("%v - %v", slot.Start.Local().Format("15:04"), slot.End.Local().Format("15:04"))))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" disabled></div></div><div class=\"three fields\"><div class=\"field\"><label>")
			if err != nil {
				return err
			}
			var_18 := `First Name`
			_, err = templBuffer.WriteString(var_18)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input name=\"firstName\"></div><div class=\"field\"><label>")
			if err != nil {
				return err
			}
			var_19 := `Last Name`
			_, err = templBuffer.WriteString(var_19)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input name=\"lastName\"></div><div class=\"field\"><label>")
			if err != nil {
				return err
			}
			var_20 := `Email`
			_, err = templBuffer.WriteString(var_20)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input name=\"email\"></div></div><button class=\"ui primary button\">")
			if err != nil {
				return err
			}
			var_21 := `Confirm`
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></form></div>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = templates.Layout().Render(templ.WithChildren(ctx, var_12), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
