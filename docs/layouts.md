GUI Layouts
===========

In this document I wish to explore the various *layouts* that are supported in
modern GUI libraries and programs, as well as the ones that are fit to be
included in hanoi.

A layout is a specification of the area in which a certain UI component can be
drawn in. They are very important in automating the correct placement of
components in their respective containers, without manually adjusting them as
the available container area also adjusts.

I am not much of a UI guy, so someone who knows more about design can fill this
in.

Centered
--------

Very important layout to consider. Used in many applications.

Properties:

- Size: User-defined
- Position: Program-defined

Fit To Parent
-------------

Will usually be used in split panes and such, when the content should always
stretch out to the full size of its container area.

Properties:

- Size: Program-defined
- Position: Beginning of drawable area

The following two are only planned, might be scrapped.

Absolute
--------

To be discouraged, but should be offered as an option

Properties:

- Size: User-defined
- Position: User-defined

Please note that the size may be defined to be greater than the size of the
container (which might be dynamically sized), or the position might be set to
somewhere outside the drawable area. This could render part of all of the
element invisible to the user.

Relative
--------

Something I've always wanted in other GUI libraries, it's a more general case of
the *general* layout, in which you're allowed to specify the percentage of the
drawable area that your element will occupy.

Properties:

- Size: Hybrid
- Position: Hybrid

Issues
======

Absolute layouts can lead to bad design very easily. A new programmer will
inevitably try to use an absolute layout because it's more liberating and lets
you put things _exactly_ where you want them, but in reality, a new programmer
will not know how to handle edge-cases and will naively hardcode some numbers
for the positioning of an element.

I wish to eliminate this with hanoi, at least for the users who will use hanoi
(read: nobody/maybe me). By only using a centered layout for every container,
and allowing for more precise positioning by the usage of more containers, a
more robust UI can be built. For example, you won't ever need to position 3
buttons in specific positions in a container for a dialog box, you'll just use a
buttonbar container, which takes the buttons as parameters, and places them in a
default, scalable way. The container itself will spawn the buttons in the view,
and give them specific drawable areas such that their fit-to-parent layout will
make them appear in the right places.

So at least for the beginning of hanoi, this is the design decision that I'm
going with. Only fit-to-parent and centered layouts, and a lot of different
containers in order to allow for more precise controls.
