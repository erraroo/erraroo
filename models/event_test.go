package models

import "testing"

func TestPopulateStackContextWithSourceMap(t *testing.T) {
	//e := Event{}
	//e.Kind ="js.error"
	//e.Payload = `{"trace":{"stack":[{"url":"http://erraroo.com/assets/erraroo-50c293f8d38851c3fa2c04bb62fb8e3f.js","func":"handled","line":536,"column":15,"context":null}]}}`

	//resources := &resourcesStore{&mockResourceGetter{}}
	//err := e.PopulateStackContext(resources)
	//assert.Nil(t, err)

	// sample data broken..
	//js, _ := e.unmarshalJSEvent()
	//frame := js.Trace.Stack[0]
	//assert.Equal(t, `      signin: function signin(user) {`, frame.PreContext[0])
	//assert.Equal(t, `        this.controllerFor('session').signin(user);`, frame.PreContext[1])
	//assert.Equal(t, `      },`, frame.PreContext[2])
	//assert.Equal(t, ``, frame.PreContext[3])
	//assert.Equal(t, `      signout: function signout() {`, frame.PreContext[4])
	//assert.Equal(t, `        this.controllerFor('session').signout();`, frame.PreContext[5])
	//assert.Equal(t, `        this.transitionTo('login');`, frame.PreContext[6])
	//assert.Equal(t, `      },`, frame.PreContext[7])
	//assert.Equal(t, ``, frame.PreContext[8])
	//assert.Equal(t, `      handled: function handled() {`, frame.PreContext[9])
	//assert.Equal(t, `        throw new Event('i threw an error');`, frame.ContextLine)
	//assert.Equal(t, `      } }`, frame.PostContext[0])
}
