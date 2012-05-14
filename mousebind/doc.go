/*
Package mousebind provides an easy to use interface to assign callback functions
to human readable button sequences.

Namely, the mousebind package exports two function types: ButtonPressFun and 
ButtonReleaseFun. Values of these types are functions, and have a method called 
'Connect' that attaches an event handler to be run when a particular button
press is issued.

This is virtually identical to the way calbacks are attached using the xevent
package, but the Connect method in the mousebind package has a couple extra
parameters that are specific to mouse bindings. Namely, the button sequence to
respond to (which is a combination of zero or more modifiers and exactly one
button), whether to establish a passive grab and whether to make the grab
synchronous or not. One can still attach callbacks to Button{Press,Release} 
events using xevent, but it will be run for *all* Button{Press,Release} events. 
(This is typically what one might do when setting up an active grab.)

Button sequence format

Button sequences are human readable strings made up of zero or more modifiers 
and exactly one button. Namely:

	[Mod[-Mod[...]]-]BUTTONNUMBER

Where 'Mod' can be one of: shift, lock, control, mod1, mod2, mod3, mod4, mod5, 
button1, button2, button3, button4, button5 or any. You can view which keys 
activate each modifier using the 'xmodmap' program. (If you don't have 
'xmodmap', you could also run the 'xmodmap' example in the examples directory.)
The 'button[1-5]' modifiers correspond to the button number in the name. (This 
implies that buttons 1 through 5 can act as both a button number and a 
modifier.)

BUTTONNUMER must correspond to a valid button number on your mouse. The best 
way to determine the numbers of each button on your mouse is to launch the xev 
program in a terminal, click the corresponding button in the new window that 
opens, and read the event output in the terminal that launched xev. Usually a
left click is button 1, a right click is button 3 and a middle click is button 
2.

An example button sequence might look like 'Mod4-Control-Shift-1'. The 
mouse binding for that button sequence is activated when all three 
modifiers---mod4, control and shift---are pressed along with the '1' button on 
your mouse.

When to issue a passive grab

One of the parameters of the 'Connect' method is whether to issue a passive 
grab or not. A passive grab is useful when you need to respond to a button press
on some parent window (like the root window) without actually focusing that 
window. Not using a passive grab is useful when you only need to read button
presses when the window is focused.

For more information on the semantics of passive grabs, please see
http://tronche.com/gui/x/xlib/input/XGrabButton.html.

Also, by default, when issuing a grab on a particular (modifiers, button) 
tuple, several grabs are actually made. In particular, for each grab requested, 
another grab is made with the "num lock" mask, another grab is made with the
"caps lock" mask, and another grab is made with both the "num lock" and "caps 
locks" masks. This allows button events to be reported regardless of whether
caps lock or num lock is enabled.

The extra masks added can be modified by changing the xevent.IgnoreMods slice. 
If you modify xevent.IgnoreMods, it should be modified once on program startup 
(i.e., before any key or mouse bindings are established) and never modified 
again.

When to use a synchronous binding

In the vast majority of cases, 'sync' in the 'Connect' method should be set to 
false, which indicates that a passive grab should be asynchronous. (The value 
of sync is irrelevant if 'grab' is false.) This implies that any events 
generated by the grab are sent to the grabbing window (the second parameter of 
the 'Connect' method) and only the grabbing window.

Sometimes, however, you might want button events to cascade down the window 
tree. That is, a button press on a parent window is grabbed, but then that 
button press should be sent to any children windows. With an asynchronous grab, 
this is impossible. With a synchronous grab, however, the button event can be 
'replayed' to all child windows. For example:

	mousebind.Initialize(XUtilValue) // call once before using mousebind package
	mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			// do something when button is pressed
			// And now replay the pointer event that fired this handler to all
			// child windows. All event processing is stopped on the
			// X server until this is called.
			xproto.AllowEvents(X.Conn(), xproto.AllowReplayPointer, 0)
		}).Connect(XUtilValue, some-window-id,
			"Mod4-Control-Shift-1", true, true)

This sort of example is precisely how reparenting window managers allow one to 
click on a window and have it be activated or "raised" *and* have the button 
press sent to the client window as well.

Note that with a synchronous grab, all event processing will be halted by the X 
server until *some* call to xproto.AllowEvents is made.

Mouse bindings on the root window example

To run a particular function whenever the 'Mod4-Control-Shift-1' button 
combination is pressed (mod4 is typically the 'super' or 'windows' key, but can
vary based on your system), use something like:

	mousebind.Initialize(XUtilValue) // call once before using mousebind package
	mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			// do something when button is pressed
		}).Connect(XUtilValue, XUtilValue.RootWin(),
			"Mod4-Control-Shift-1", false, true)

Note that we issue a passive grab because Button{Press,Release} events on the 
root window will only be reported when the root window has focus if no grab 
exists.

Mouse bindings on a window you create example

This code snippet attaches an event handler to some window you've created 
without using a grab. Thus, the function will only be activated when the button 
sequence is pressed and your window has focus.

	mousebind.Initialize(XUtilValue) // call once before using mousebind package
	mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			// do something when button is pressed
		}).Connect(XUtilValue, your-window-id, "Mod4-t", false, false)

Run a function on all button press events example

This code snippet actually does *not* use the mousebind package, but illustrates
how the Button{Press,Release} event handlers in the xevent package can still be 
useful. Namely, the mousebind package discriminates among events depending upon 
the button sequences pressed, whereas the xevent package is more general: it 
can only discriminate at the event level.

	xevent.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			// do something when any button is pressed
		}).Connect(XUtilValue, your-window-id)

This is the kind of handler you might use to capture all button press events. 

More examples

A complete working example using the mousebind package can be found in
'simple-mousebinding' in the examples directory of the xgbutil package.

*/
package mousebind