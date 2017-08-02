# $NetBSD$

BUILDLINK_TREE+=	png++

.if !defined(PNG++_BUILDLINK3_MK)
PNG++_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.png+++=	png++>=0.2.9
BUILDLINK_PKGSRCDIR.png++?=	../../graphics/png++

#.include "../../graphics/png/buildlink3.mk"
.endif	# PNG++_BUILDLINK3_MK

BUILDLINK_TREE+=	-png++
