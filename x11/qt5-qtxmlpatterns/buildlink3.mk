# $NetBSD: buildlink3.mk,v 1.25 2019/04/03 00:33:20 ryoon Exp $

BUILDLINK_TREE+=	qt5-qtxmlpatterns

.if !defined(QT5_QTXMLPATTERNS_BUILDLINK3_MK)
QT5_QTXMLPATTERNS_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.qt5-qtxmlpatterns+=	qt5-qtxmlpatterns>=5.9.1
BUILDLINK_ABI_DEPENDS.qt5-qtxmlpatterns+=	qt5-qtxmlpatterns>=5.12.2nb1
BUILDLINK_PKGSRCDIR.qt5-qtxmlpatterns?=		../../x11/qt5-qtxmlpatterns

BUILDLINK_INCDIRS.qt5-qtxmlpatterns+=	qt5/include
BUILDLINK_LIBDIRS.qt5-qtxmlpatterns+=	qt5/lib
BUILDLINK_LIBDIRS.qt5-qtxmlpatterns+=	qt5/plugins

.include "../../x11/qt5-qtbase/buildlink3.mk"
.endif	# QT5_QTXMLPATTERNS_BUILDLINK3_MK

BUILDLINK_TREE+=	-qt5-qtxmlpatterns
