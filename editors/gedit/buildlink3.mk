# $NetBSD: buildlink3.mk,v 1.48 2018/11/14 22:21:23 kleink Exp $

BUILDLINK_TREE+=	gedit

.if !defined(GEDIT_BUILDLINK3_MK)
GEDIT_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.gedit+=	gedit>=2.12.1nb4<3
BUILDLINK_ABI_DEPENDS.gedit+=	gedit>=2.30.4nb32
BUILDLINK_PKGSRCDIR.gedit?=	../../editors/gedit

.include "../../x11/gtksourceview2/buildlink3.mk"
.endif # GEDIT_BUILDLINK3_MK

BUILDLINK_TREE+=	-gedit
